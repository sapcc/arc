package vendor

//This is a dummy package which is not used in arcs codebase
//it just here to make godep save ./... not discard
//dependecies that are only used on windows

import (
	_ "github.com/go-ole/go-ole"
	_ "github.com/go-ole/go-ole/oleutil"
)
