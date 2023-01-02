package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

const localPartImageStorage = "/files/parts/"
const localSetImageStorage = "/files/sets/"
const localIconStorage = "/files/buttons/"
const cssFilesRoot = "/files/css/"
const htmlFilesRoot = "/files/html/"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeView)
	router.HandleFunc("/addset", addSetView)
	router.HandleFunc("/viewset", singleSetView)
	router.HandleFunc("/createdatabase", createDatabaseView)
	router.HandleFunc("/viewimage", viewImage)
	router.HandleFunc("/geticon", getIcon)
	router.HandleFunc("/modifycount", modifyCount)
	router.HandleFunc("/addcolorlist", addColorList)
	router.HandleFunc("/ping", pingView)
	router.HandleFunc("/style", cssView)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func decodeSet(data []byte) Inventory {
	v := Inventory{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return v
}

func decodeColors(data []byte) ColorCatalog {
	v := ColorCatalog{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return v
}

func homeView(w http.ResponseWriter, r *http.Request) {
	setList, err := getSetList()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, getHtmlFileContents("setlistheader"))
	for _, item := range setList {
		fmt.Fprintf(w, getHtmlFileContents("setlistitem"), item.ID, item.ID, item.Name, item.ID, item.Name, item.Description)
	}
	fmt.Fprintf(w, getHtmlFileContents("setlistformfooter"))
}

func singleSetView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	idstr := query.Get("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Can't find set, %s", err.Error())
		return
	}

	if id > 0 {
		fmt.Fprintf(w, getHtmlFileContents("partlistheader"), id)
		setParts, err := getPartList(id)
		if err != nil {
			fmt.Fprintf(w, "Can't find set, %s", err.Error())
			return
		}
		legoColors, err := getLegoColors()
		if err != nil {
			fmt.Fprintf(w, "Can't find legocolors, %s", err.Error())
		}
		for _, part := range setParts {
			colorName, found := legoColors[part.ColorID]
			if !found {
				colorName = "ingen"
			}
			//colorName := From(legoColors).Where(func(c LegoColor) bool { return c.Number == part.ColorID }).Select(func(c LegoColor) string { return c.Name }).First()
			fmt.Fprintf(w, getHtmlFileContents("partlistitem"), part.Partnumber, part.ColorID, part.Partnumber, part.Partnumber, part.ColorID, part.Partnumber, part.Partnumber, part.ColorID, part.RequiredQty, part.Partnumber, part.ColorID, part.FoundQty, part.Partnumber, part.ColorID, part.Partnumber, part.ColorID, colorName, part.Description)
		}
	}
}

func viewImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	ID := query.Get("id")
	imageType := query.Get("type")

	var imagePath string
	if imageType == "p" {
		colorstr := query.Get("color")
		color, err := strconv.Atoi(colorstr)
		if err != nil {
			fmt.Fprintf(w, "Can't get part color, %s", err.Error())
			return
		}
		imagePath = localPartImageStorage + ID + "-" + strconv.Itoa(color) + ".png"
	} else {
		imagePath = localSetImageStorage + ID + ".png"
	}

	if _, err := os.Stat(imagePath); err == nil {
		http.ServeFile(w, r, imagePath)
	} else {
		imagePath = localSetImageStorage + "na.png"
		http.ServeFile(w, r, imagePath)
	}
}

func getIcon(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	name := query.Get("name")

	imagePath := localIconStorage + name
	http.ServeFile(w, r, imagePath)
}

func createDatabaseView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getHtmlFileContents("createdatabaseheader"))
	err := createDatabase()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, "Success!")
	}
	fmt.Fprintf(w, getHtmlFileContents("createdatabasefooter"))
}

func addSetView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, getHtmlFileContents("addset"))
	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(1000000)
		if err != nil {
			fmt.Fprintf(w, "FILE_TOO_BIG")
			return
		}

		file, _, err := r.FormFile("setfilename")
		if err != nil {
			fmt.Println("Error retrieving the file")
			fmt.Println(err)
			return
		}
		defer file.Close()

		setName := r.FormValue("setname")
		setDescription := r.FormValue("setdescription")
		setImageURL := r.FormValue("setimageurl")

		var legoset LegoSet
		legoset.Name = setName
		legoset.Description = setDescription
		setID := saveSet(legoset)
		saveImageByURL(setID, setImageURL)

		// decode file
		content, err := ioutil.ReadAll(file)
		inventory := decodeSet(content)
		// save items
		var legoParts []LegoPart
		for _, item := range inventory.Items {
			// retrieve part info: description and image
			description, _ := getPartImageAndDescription(item.ItemID, item.Color)
			var legoPart LegoPart
			legoPart.Description = description
			legoPart.ColorID = item.Color
			legoPart.Partnumber = item.ItemID
			legoPart.SetID = setID
			legoPart.RequiredQty = item.MinQty
			legoPart.FoundQty = 0
			legoParts = append(legoParts, legoPart)
		}
		saveParts(legoParts)
	}
}

func modifyCount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	setID := r.FormValue("setid")
	id := r.FormValue("id")

	re := regexp.MustCompile(`(.*)-(.*)-(.*)`)
	subMatches := re.FindStringSubmatch(id)

	if len(subMatches) == 4 {
		direction := subMatches[1]
		partNumber := subMatches[2]
		colorID := subMatches[3]
		foundqty, err := setPartFoundQuantity(setID, partNumber, colorID, direction)
		if err != nil {
			w.Write([]byte("-1"))
			return
		}
		w.Write([]byte(strconv.Itoa(foundqty)))
		return
	}

	// her bør den nye foundqty nok returneres...
	w.Write([]byte("-1"))
}

func addColorList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, getHtmlFileContents("setcolorlist"))
	} else if r.Method == "POST" {
		err := r.ParseMultipartForm(1000000)
		if err != nil {
			fmt.Fprintf(w, "FILE_TOO_BIG")
			return
		}

		file, _, err := r.FormFile("colorfilename")
		if err != nil {
			fmt.Println("Error retrieving the file")
			fmt.Println(err)
			return
		}
		defer file.Close()

		// decode file
		content, err := ioutil.ReadAll(file)
		colors := decodeColors(content)
		// save items
		var legoColors []LegoColor
		for _, item := range colors.Items {
			var legoColor LegoColor
			legoColor.Number = item.Color
			legoColor.Name = item.ColorName
			legoColors = append(legoColors, legoColor)
		}
		saveColors(legoColors)
		fmt.Fprintf(w, "done")
	}
	fmt.Fprintf(w, "something happened...")
}

func cssView(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	cssFile := cssFilesRoot + filename + ".css"

	if _, err := os.Stat(imagePath); err == nil {
		http.ServeFile(w, r, cssFile)
	} else {
		cssFile = cssFileRoot + "empty.css"
		http.ServeFile(w, r, cssFile)
	}
}

func getHtmlFileContents(filename string) string {
    file, err := os.Open(htmlFilesRoot + filename + ".html")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    b, err := ioutil.ReadAll(file)
    return string(b)
}

func pingView(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}
