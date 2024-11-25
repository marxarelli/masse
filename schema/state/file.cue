package state

#File: {
	file!:    #FileAction | [#FileAction, ...#FileAction]
	options?: #Option | [#Option, ...#Option]
}

#FileAction: {
	#Copy | #Mkfile | #Mkdir | #Rm
}
