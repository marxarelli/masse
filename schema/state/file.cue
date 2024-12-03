package state

#File: {
	file!:    #FileAction | [#FileAction, ...#FileAction]
	options?: #Constraint | [#Constraint, ...#Constraint]
}

#FileAction: {
	#Copy | #Mkfile | #Mkdir | #Rm
}
