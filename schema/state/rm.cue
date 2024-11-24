package state

#Rm: {
	rm!:      string
	options?: #RmOption | [#RmOption, ...#RmOption]
}

#RmOption: {
	#AllowNotFound | #AllowWildcard
}

#AllowNotFound: {
	allowNotFound!: bool
}

#AllowWildcard: {
	allowWildcard!: bool
}
