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
	"strings"

	"github.com/futtie/LegoCollector/database"
	"github.com/futtie/LegoCollector/structs"

	"github.com/gorilla/mux"
)

const localPartImageStorage = "./files/parts/"
const localSetImageStorage = "./files/sets/"
const localIconStorage = "./files/buttons/"
const cssFilesRoot = "./files/css/"
const htmlFilesRoot = "./files/html/"

var db	 *database.Database 

func main() {
	//db = database.InitDatabase("legouser:legopassword@tcp(localhost:3306)/legoparts")
	db = database.InitDatabase("legouser:legopassword@tcp(dbserver:3306)/legoparts")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeView)
	router.HandleFunc("/addset", addSetView)
	router.HandleFunc("/getmissinglist", missingListView)
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

func decodeSet(data []byte) structs.Inventory {
	v := structs.Inventory{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return v
}

func encodeSet(data structs.Inventory) string {
	b, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return string(b)
}

func decodeColors(data []byte) structs.ColorCatalog {
	v := structs.ColorCatalog{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return v
}

func homeView(w http.ResponseWriter, r *http.Request) {
	setList, err := db.GetSetList()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, getHtmlFileContents("setlistheader"))
	for _, item := range setList {
		fmt.Fprintf(w, getHtmlFileContents("setlistitem"), item.ID, item.ID, item.Name, item.ID, item.Name, item.Description, item.RequiredCount, item.FoundCount, item.ID)
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
		setParts, err := db.GetPartList(id)
		if err != nil {
			fmt.Fprintf(w, "Can't find set, %s", err.Error())
			return
		}
		legoColors, err := db.GetLegoColors()
		if err != nil {
			fmt.Fprintf(w, "Can't find legocolors, %s", err.Error())
		}
		for _, part := range setParts {
			colorName, found := legoColors[part.ColorID]
			if !found {
				colorName = "ingen"
			}
			
			trclass := "class"
			if part.RequiredQty == part.FoundQty {
				trclass = "class=\"allpartsfound\""
			}
			//colorName := From(legoColors).Where(func(c LegoColor) bool { return c.Number == part.ColorID }).Select(func(c LegoColor) string { return c.Name }).First()
			fmt.Fprintf(w, getHtmlFileContents("partlistitem"), part.Partnumber, part.ColorID, trclass, part.Partnumber, part.Partnumber, part.ColorID, part.Partnumber, part.Partnumber, part.ColorID, part.RequiredQty, part.Partnumber, part.ColorID, part.FoundQty, part.Partnumber, part.ColorID, part.Partnumber, part.ColorID, part.Partnumber, part.ColorID, colorName, part.Description)
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
	err := db.CreateDatabase()
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

		var legoset structs.LegoSet
		legoset.Name = setName
		legoset.Description = setDescription
		setID := db.SaveSet(legoset)
		saveImageByURL(setID, setImageURL)

		// decode file
		content, err := ioutil.ReadAll(file)
		inventory := decodeSet(content)
		// save items
		var legoParts []structs.LegoPart
		for _, item := range inventory.Items {
			// retrieve part info: description and image
			description, _ := getPartImageAndDescription(item.ItemID, item.Color)
			var legoPart structs.LegoPart
			legoPart.Description = description
			legoPart.ColorID = item.Color
			legoPart.Partnumber = item.ItemID
			legoPart.SetID = setID
			legoPart.RequiredQty = item.MinQty
			legoPart.FoundQty = 0
			legoParts = append(legoParts, legoPart)
		}
		db.SaveParts(legoParts)
	}
	homeView(w, r)
}

func missingListView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	setID := r.FormValue("setid")
	id, err := strconv.Atoi(setID)
	if err != nil {
		fmt.Fprintf(w, "Can't find set, %s", err.Error())
		return
	}
	setParts, err := db.GetPartList(id)
	if err != nil {
		fmt.Fprintf(w, "Could not get parts list")
		return
	}
	var inventory structs.Inventory

	for _, part := range setParts {
		var item structs.Item
		item.ItemType = "P"
		item.ItemID = part.Partnumber
		item.Color = part.ColorID
		item.Maxprice = -1
		item.MinQty = part.RequiredQty - part.FoundQty
		item.Condition = "U"
		item.Notify = "N"
		inventory.Items = append(inventory.Items, item)
	}
	xmlfile := encodeSet(inventory)
	xmlfile = strings.Replace(xmlfile, "<Inventory>", "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n<INVENTORY xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">", -1)
	xmlfile = strings.Replace(xmlfile, "</Inventory>", "</INVENTORY>", -1)
	fmt.Fprintf(w, xmlfile)
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
		foundqty, err := db.SetPartFoundQuantity(setID, partNumber, colorID, direction)
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
		var legoColors []structs.LegoColor
		for _, item := range colors.Items {
			var legoColor structs.LegoColor
			legoColor.Number = item.Color
			legoColor.Name = item.ColorName
			legoColors = append(legoColors, legoColor)
		}
		db.SaveColors(legoColors)
		fmt.Fprintf(w, "done")
	}
	homeView(w, r)
}

func cssView(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	cssFile := cssFilesRoot + filename + ".css"
	http.ServeFile(w, r, cssFile)
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
