package main

// HTMLSetListHeader ...
const HTMLSetListHeader = `<html>
<head>
	<title>Setlist</title>
</head>	
<body>
	<form enctype="multipart/form-data" action="/addset" method="post">
		<table>
			<tbody>
`

// HTMLSetListItemHeader ...
const HTMLSetListItemHeader = `
				<tr>
					<td>Set</td>
					<td>Image</td>
					<td>Description</td>
				</tr>
`

// HTMLSetListitem ...
const HTMLSetListitem = `
				<tr id="%d">
					<td><a href="/viewset?id=%d">%s</a></td>
					<td>%s</td>
					<td>%s</td>
				</tr>
`

// HTMLSetListFooter ...
const HTMLSetListFooter = `
				</tbody>
		</table>
	</form>
</body>
</html>
`

// HTMLPartListHeader ...
const HTMLPartListHeader = `<html>
<head>
	<title>Partlist</title>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
	<script>
		var setID = %d;
		function postdata(item){
			var imageid = item.id
			var partnumbercolor = imageid.substring(imageid.indexOf('-'));
			var foundtag = '#found'.concat(partnumbercolor)
			$.ajax({
				url: "/modifycount",
				type: "POST",
				data: { 
					id: imageid,
					setid: setID
				},
				success : function(data) {
					$(foundtag).text(data);
					var reqtag = '#req'.concat(partnumbercolor)
					var reqcount = $(reqtag).text()
					if (data === reqcount)
						$('#tr'.concat(partnumbercolor)).hide();
				},
			});
		};
	</script>
</head>	
<body>
	<form enctype="multipart/form-data" action="/addset" method="post">
		<table>
			<colgroup>
				<col style="background-color:#efefef">
				<col style="background-color:white">
				<col style="background-color:#efefef">
				<col style="background-color:white">
				<col style="background-color:#efefef">
				<col style="background-color:white">
				<col style="background-color:#efefef">
	  		</colgroup>
			<tbody>
`

// HTMLPartListItemHeader ...
const HTMLPartListItemHeader = `
				<tr>
					<td>Del</td>
					<td width="15%%">Billede</td>
					<td>Antal nødvendig</td>
					<td>Antal fundet</td>
					<td>Fortryd / Fundet</td>
					<td>Farve</td>
					<td>Beskrivelse</td>
				</tr>
`

// HTMLPartListItem ...
const HTMLPartListItem = `
				<tr id="tr-%s-%d">
					<td>%s</td>
					<td><img src="/viewimage?id=%s&color=%d" alt="%s" width="90%%"></td>
					<td><h3 id="req-%s-%d" align=center>%d</h3></td>
					<td><h3 id="found-%s-%d" align=center>%d</h3></td>
					<td>
						<img id="down-%s-%d" onclick="postdata(this)" src="/geticon?name=down1.png" width="50px" height="50px"/>
						<img id="up-%s-%d" onclick="postdata(this)" src="/geticon?name=up1.png" width="50px" height="50px"/>
					</td>
					<td>%s</td>
					<td>%s</td>
				</tr>
`

// HTMLPartListFooter ...
const HTMLPartListFooter = `
				</tbody>
		</table>
	</form>
</body>
</html>
`

// HTMLAddSet ...
const HTMLAddSet = `<html>
	<head>
		<title>Add set</title>
	</head>	
	<body>
		<form enctype="multipart/form-data" action="/addset" method="post">
			<table>
				<tbody>
					<tr>
						<td>Sæt navn eller nummer</td>
						<td><input type="text" name="setname"></td>
					</tr>
					<tr>
						<td>Sæt beskrivelse</td>
						<td><input type="text" name="setdescription"></td>
					</tr>
					<tr>
						<td>Sæt billede</td>
						<td><input type="text" name="setimageurl"></td>
					</tr>
					<tr>
						<td>Inventory XML fra Bricklink fra wanting liste</td>
						<td><input type="file" name="setfilename"></td>
					</tr>
					<tr>
						<td><input type="submit" value="Tilbage">
						<td><input type="reset" value="Nulstil"><input type="submit" value="Gem">
					</tr>
				</tbody>
			</table>
		</form>
	</body>
</html>
`

// HTMLCreateDatabaseHeader ...
const HTMLCreateDatabaseHeader = `<html>
	<head>
		<title>Create database</title>
	</head>
	<body>
`

// HTMLCreateDatabaseFooter ...
const HTMLCreateDatabaseFooter = ` 
	</body>
</html>
`

// HTMLSetColorList ...
const HTMLSetColorList = `<html>
	<head>
		<title>Set color list</title>
	</head>	
	<body>
		<form enctype="multipart/form-data" action="/addcolorlist" method="post">
			<table>
				<tbody>
					<tr>
						<td>Color list XML fra Bricklink</td>
						<td><input type="file" name="colorfilename"></td>
					</tr>
					<tr>
						<td><input type="submit" value="Tilbage">
						<td><input type="reset" value="Nulstil"><input type="submit" value="Gem">
					</tr>
				</tbody>
			</table>
		</form>
	</body>
</html>
`
