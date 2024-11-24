package state

#File: {
	file!:    #FileAction | [#FileAction, ...]
	options?: #Option | [#Option, ...]
}

#FileAction: {
	#Copy | #Mkfile | #Mkdir | #Rm
}
