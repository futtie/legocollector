# LegoCollector
Collect all parts for a lego set

v indlæse en ny model samt billede af model
v xml/html decoder
v database handler
v visning af siden
v visning af enkelt model med nedskrivning
v dan billede for del
v gem billede for model


TODO: 
v hent billede ud fra partnumber/color
v generer billedeurl og hent/gem billede hvis det ikke findes
v dan og vis side med sætliste
v dan og vis side med enkelt sæt
v gem sæt og parts og gem tilhørende billeder
v mulighed for at vise enten
    v fuld liste med disabled rækker hvor alle dele er fundet
    v nettoliste uden de rækker hvor alle dele er fundet
- mulighed for genaktivering af dele - Won't Do
- nulstil et helt sæt inkl. "er du sikker" 

v Skift skrifttype ffs

v skift baggrundsfarve på de rækker hvor alle er fundet
  v skift tilbage til normal ved fortryd og overhold de farvevalg der er truffet for kolonnerne
  v lige nu fungerer skift til fundet ved onclick, ikke ved load - det skal muligvis klares ved at tilføje en klasse når alt er fundet og fjerne den igen ved fortryd
  
v 0.3
  v rette fejl

v 0.4
  v rette databasefejl med for mange forbindelser

v 0.5
  v Gør muligt at tilføje 10 ad gangen (husk at der ikke må blive flere end nødvendigt)
  v sorter efter farve

v 0.6
  v Større indtastningsfelter (input)
  v "Sæt billede" skal have en "Choose file"
  v "/addset" skal returnere til "/setlist" eller "/partlist"
  v mulighed for at gemme de manglende stumper som XML, som kan importeres i Bricklink
  v større knapper
  v sætliste skal vise "totalt fundet antal" og "totalt antal"

1.0
- mulighed for zoom på siden
