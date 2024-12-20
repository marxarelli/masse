package state

#Rm: {
	rm!:      string
	options?: #RmOption | [#RmOption, ...#RmOption]
}

#RmOption: {
	#AllowNotFound | #Wildcard
}
