\input header

@** Introduction.
This is an implementation of \.{ahist} command for \.{Acme}. It tracks all search requests in \.{Acme}'s window to a separate window.


@** Implementation.
@c

@i license

package main

import (
	@<Imports@>
)@#

var (
	@<Global variables@>
)@#

type (
	@<Types@>
)@#

@* Startup.
@c
func main () {
	@<Store a name of the program@>
	@<Obtaining of |id| of a window@>
	@<Open window |w| by |id|@>
	@<Change the name of the program in the tag@>
	@<Read |name| of the window@>
	@<Start history processing@>
	@<Processing window events@>
}

@
@<Imports@>=
"fmt"
"os"

@ Let's define |dbg| flag and will switch it by |ahist+| and |ahist-|.
@<Global variables@>=
dbg bool

@
@<Switch debug output on@>=
dbg=true
debug("debug has been switched on\n")

@
@<Switch debug output off@>=
debug("debug has been switched off\n")
dbg=false

@
@c
func debug(f string, args...interface{}) {
	if dbg {
		fmt.Fprintf(os.Stderr, f, args...)
	}
}

@
@<Imports@>=
"strconv"

@
@<Global variables@>=
id int

@
@<Obtaining of |id| of a window@>=
{
	var err error
	id, err=strconv.Atoi(os.Getenv("winid"))
	if err!=nil {
 	       return
	}
}

@
@<Global variables@>=
tagname string

@
@<Store a name of the program@>=
tagname=os.Args[0]
if n:=strings.LastIndex(tagname, "/"); n!=-1 {
	tagname=tagname[n:]
}
debug("tagname:%s\n", tagname)

@ We change \.{ahist} into \.{-ahist} to add a possibility to switch \.{ahist} off.
@<Change the name of the program in the tag@>=
{
	del:=[]string{tagname, "-"+tagname, "-"+tagname+"+", "-"+tagname+"-"}
	add:=[]string{"-"+tagname}
	changeTag(w, del, add)
}

@ On exit we should make an opposite change.
@<Cleanup@>=
{
	del:=[]string{tagname, "-"+tagname, "-"+tagname+"+", "-"+tagname+"-"}
	add:=[]string{tagname}
	changeTag(w, del, add)
}


@* Events handling.

@
@<Imports@>=
"github.com/santucco/goacme"

@
@<Processing window events@>=
@<Fix tag of the window@>
for {
	ev, err:=w.ReadEvent()
	if err!=nil {
		@<Cleanup@>
		return
	}
	@<Process main window@>
}

@
@<Open window |w| by |id|@>=
w, err:=goacme.Open(id)
if err!=nil {
	debug("cannot open a window with id %d: %s\n", id, err)
	return
}
defer w.Close()

@
@<Imports@>=
"strings"

@
@<Global variables@>=
name string

@
@<Read |name| of the window@>=
{
	f, err:=w.File("tag")
	if err!=nil {
		debug("cannot read from 'tag' of the window with id %d: %s\n", id, err)
		return
	}
	if _, err:=f.Seek(0, 0); err!=nil {
		debug("cannot seek to the start 'tag' of the window with id %d: %s\n", id, err)
		return
	}
	var b [1000]byte
	n, err:=f.Read(b[:])
	if err!=nil {
		debug("cannot read tag of the window with id %d: %s\n", id, err)
		return
	}
	ss:=strings.Split(string(b[:n]), " ")
	if len(ss)==0 {
		return
	}
	name=string(ss[0])
}

@
@<Process main window@>=
@<Process and continue if it is not |Look| in any form@>
@<Process |Look|@>
@<Read addr into |b, e|@>
@<Show dot@>
@<Write history@>

@ |b,e| address pair is taken from the |ev| event.
@<Process and continue if it is not |Look| in any form@>=
debug("incoming event: %+v\n", ev)
s:=""
b:=ev.Begin
e:=ev.End
type_switch: switch {
	case ev.Type==goacme.Look|goacme.Tag:
		@<Process in case of a request by |B3| mouse button in the tag@>
	case ev.Type==goacme.Look:
		@<Process in case of a request by |B3| command in the body@>
	case ev.Type==goacme.Execute || ev.Type==goacme.Execute|goacme.Tag:
		@<Process in case of executing a command in the body or tag@>
	case ev.Type==goacme.Insert || ev.Type==goacme.Delete:
		@<Fix tag of the window@>
		continue
	default:
		@<Unread event and continue@>
}

@ We take a search string from |ev| event and set dot.
Also we have to clean |b,e| because it is an address in the tag.
@<Process in case of a request by |B3| mouse button in the tag@>=
b,e=0,0
s=ev.Text
if len(ev.Arg)>0 {
	s+=" "+ev.Arg
}
@<Set addr to dot@>

@ We take a search string and address from |ev| event.
@<Process in case of a request by |B3| command in the body@>=
s=ev.Text
if len(ev.Arg)>0 {
	s+=" "+ev.Arg
}
@<Set addr to |b, e|@>

@ For |Look| command we set address and continue processing.
|ahist| command we just ignore to avoid duplicates.
|-ahist| command makes cleanups and processes to exit.
|ahist+| and |ahist-| switch debug output on and off.
All other commands are written back to |"event"| file and |fallthrough|
to the next case, where a status of the window is checked.
@<Process in case of executing a command in the body or tag@>=
switch strings.TrimSpace(ev.Text) {
	case "Look":
		@<Process in case of executing |Look| command@>
		break type_switch
	case tagname:
		continue
	case "-"+tagname+"+":
		fallthrough
	case "-"+tagname+"-":
		fallthrough
	case "-"+tagname:
		debug("exiting\n")
		@<Cleanup@>
		return
	case tagname+"+":
		@<Switch debug output on@>
		continue
	case tagname+"-":
		@<Switch debug output off@>
		continue
}
w.UnreadEvent(ev)
fallthrough

@ We take a search string from an argument of |Look| command.
Current address is set to dot, then |b,e| pair is set to the current address.
@<Process in case of executing |Look| command@>=
s=ev.Arg
@<Set addr to dot@>
@<Read addr into |b, e|@>

@
@<Unread event and continue@>=
w.UnreadEvent(ev)
continue

@ If the |ev| event contains a search string, use it.
Otherwise we should read selected the string from the window's body and read its address into |b,e|.
@<Process |Look|@>=
{
	if len(s)>0 {
		@<Make a search of |s|@>
	} else {
		@<Look for selected string@>
		@<Read addr into |b, e|@>
	}
}

@
@<Look for selected string@>={
	@<Read selected string from |"xdata"| file to |s|@>
	@<Make a search of |s|@>
}

@
@<Set addr to dot@>=
if w.WriteCtl("addr=dot")!=nil {
	@<Unread event and continue@>
}
debug("set addr to dot\n")

@
@<Read selected string from |"xdata"| file to |s|@>={
	d, err:=w.File("xdata")
	if err!=nil {
		debug("cannot read from 'xdata' of the window with id %d: %s\n", id, err)
		@<Unread event and continue@>
	}

	buf:=make([]byte, e-b+1)

	for n, _:=d.Read(buf); n>0; n, _=d.Read(buf) {
		s+=string(buf[:n])
	}
	debug("read address from xdata b: %v, e: %v\n", b, e)
}

@
@<Set addr to |b, e|@>=
if err:=w.WriteAddr("#%d,#%d", b, e); err!=nil {
	debug("cannot write to 'addr' of the window with id %d: %s\n", id, err)
	@<Unread event and continue@>
}
debug("set addr to %d, %d\n", b, e)

@ We need to story previous history |entry| for the case, when |Look| in a tag is executed
but without selected text. In the case a search string is taken from \.{Acme}.
We take it from |lentr|
@<Global variables@>=
lentr entry

@ Let's add |empty| function for |entry|
@c
func (this entry) empty() bool {
	return this.b==this.e
}

@ Search is processed by writing |"/<regex>/"| to |"addr"| file,
but before regex-specific symbols of |s| have to be escaped
In the case of empty search string we take it from |lentr|.
Also we write the current position with the string to the history to track the search,
because it already has a place.
@<Make a search of |s|@>={
	debug("last entry : %v\n", lentr)
	if len(s)==0 {
		if !lentr.empty() {
			b=lentr.b
			e=lentr.e
			s=lentr.s
			@<Set addr to |b, e|@>
		}
	} else if b!=e {
		lentr=entry{b,e,s}
		@<Write history@>
	}
	es:=escapeSymbols(s)
	debug("escaped search string: %q\n", es)
	if err:=w.WriteAddr("/%s/", es); err!=nil {
		debug("cannot write to 'addr' of the window with id %d: %s\n", id, err)
		@<Unread event and continue@>
	}
}

@
@<Read addr into |b, e|@>=
b, e, err=w.ReadAddr()
if err!=nil {
	@<Unread event and continue@>
}
debug("read address b: %v, e: %v\n", b, e)

@
@<Set dot to addr@>=
if w.WriteCtl("dot=addr")!=nil {
	debug("cannot write to 'ctl' of the window with id %d: %s\n", id, err)
	@<Unread event and continue@>
}
debug("set dot to addr\n")

@
@<Show dot@>=
@<Set dot to addr@>
if w.WriteCtl("show")!=nil {
	debug("cannot write to 'ctl' of the window with id %d: %s\n", id, err)
	@<Unread event and continue@>
}
debug("show dot\n")

@ \.{Acme} does not produce standard commands in case of opened |"event"| file.
So we have to add command |"Put"| in case of the window is modified and |"Undo"| and |"Redo"| commands too.
@<Fix tag of the window@>=
{
	_, _, _, _, d, _, _, _, err:=w.ReadCtl()
	if err!=nil {
		debug("cannot read from 'ctl' of the window with id %d: %s\n", id, err)
	} else {
		debug("dirty: %v\n", d)
		del:=[]string{"Put", "Undo", "Redo"}
		var add []string
		if d {
			add=append(add, "Put")
		}
		add=append(add, "Undo", "Redo")
		changeTag(w, del, add)
	}
}

@ Removing added commands on exit
@<Cleanup@>=
{
	del:=append([]string{}, "Put", "Undo", "Redo")
	changeTag(w, del, nil)
}

@* Tracking search requests .

We create a window with history of search requests and make separated goroutine to process events from the window.

@
@<Types@>=
entry struct {
	b, e int
	s string
}

@ Special |histch| channel is received |entry| to print them in the window
@<Global variables@>=
histch chan entry=make(chan entry)

@ On exit we should signal the goroutine to stop processing. It is made by closing |histch| channel
@<Cleanup@>=
close(histch)

@
@<Variables outside the loop@>=
var hch <-chan *goacme.Event

@ The goroutine handles two variants of events.
@<Start history processing@>=
go func(){
	@<Variables outside the loop@>
	for {
		select {
			case entr, ok:=<-histch:
				@<Process |entr| entry from |histch|@>
			case ev, ok:=<-hch:
				@<Process |ev| event from |hch| event channel of the window@>
		}
	}
}()

@
@<Variables outside the loop@>=
var h *goacme.Window

@ Events from |histch| channel is written to the history.
Before writing a history entry we look for the address in the history window and
write the entry only if it has not been found.
@<Process |entr| entry from |histch|@>=
if !ok {
	if h!=nil {
		h.Del(true)
		h.Close()
		h=nil
	}
	return
}
@<Open history window, if it does not exist@>
if h.WriteAddr("/#%d,#%d/", entr.b, entr.e)!=nil {
	debug("writing to the history %d,%d\n", entr.b, entr.e)
	h.Write([]byte(fmt.Sprintf("%s:#%d,#%d %q\n", name, entr.b, entr.e, entr.s)))
	h.WriteCtl("clean")
}
debug("selecting the current position #%d,#%d in the history\n", entr.b, entr.e)
es:=fmt.Sprintf("#%d,#%d", entr.b, entr.e)
@<Make a selection of the current search request@>

@ Event from |hch| channel is checked for a case the channel is close.
In the case that means the history window is closed and we clear |h| and |hch|.
Otherwise we just write the event back.
@<Process |ev| event from |hch| event channel of the window@>=
if !ok {
	debug("history is closed\n")
	h.Del(true)
	h.Close()
	h=nil
	hch=nil
	continue
}
h.UnreadEvent(ev)
if ev.Type==goacme.Look {
	debug("incoming event: %+v\n", ev)
	debug("selecting the current position %q in the history\n", ev.Text)
	es:=escapeSymbols(ev.Text)
	@<Make a selection of the current search request@>
}

@
@<Make a selection of the current search request@>=
if err:=h.WriteAddr("/%s/-+", es); err!=nil {
	debug("writing of addr failed: %s\n", err)
} else if err:=h.WriteCtl("dot=addr\nshow"); err!=nil {
	debug("writing of ctl failed: %s\n", err)
}

@
@<Write history@>=
debug("request to store a history: %v,%v %q\n", b, e, s)
histch<-entry{b:b, e:e, s:s}

@ If the history window |h| does not exist, we create it.
@<Open history window, if it does not exist@>=
if h==nil {
	var err error
	if h, err=goacme.New(); err!=nil {
		return
	}
	h.WriteCtl("name %s", name+"+History")
	if hch, err=h.EventChannel(1, goacme.AllTypes); err!=nil {
		return
	}
}

@ |changeTag| function.

We read the tag of |w| window, remove all commands from |del| list 
and add all commands from |add| list.
@c
func changeTag(w *goacme.Window, del[]string, add []string) {
	if add==nil && del==nil {
		return
	}
	@<Read a tag of |w| into |s|@>
	@<Split tag into |tag| fields after the pipe symbol@>
	@<Compose |newtag|@>
	@<Clear the tag and write |newtag| to the tag@>
}

@
@<Read a tag of |w| into |s|@>=
f, err:=w.File("tag")
if err!=nil {
	debug("cannot read from 'tag' of the window with id %d: %s\n", id, err)
	return
}
if _, err:=f.Seek(0, 0); err!=nil {
	debug("cannot seek to the start 'tag' of the window with id %d: %s\n", id, err)
	return
}
var b [1000]byte
n, err:=f.Read(b[:])
if err!=nil {
	debug("cannot read tag of the window with id %d: %s\n", id, err)
	return
}
s:=string(b[:n])

@
@<Split tag into |tag| fields after the pipe symbol@>=
if n=strings.LastIndex(s, "|"); n==-1 {
	n=0
} else {
	n++
}
s=s[n:]
s=strings.TrimLeft(s, " ")
tag:=strings.Split(s, " ")

@
@<Compose |newtag|@>=
newtag:=append([]string{}, "")
@<Every part is contained in |del| we remove from |tag|@>
newtag=append(newtag, add...)
newtag=append(newtag, tag...)

@
@<Every part is contained in |del| we remove from |tag|@>=
for _, v:=range del {
	for i:=0; i<len(tag); {
		if tag[i]!=v {
			i++
			continue
		}
		copy(tag[i:], tag[i+1:])
		tag=tag[:len(tag)-1]
	}
}

@
@<Clear the tag and write |newtag| to the tag@>=
s=strings.Join(newtag, " ")	
if err:=w.WriteCtl("cleartag"); err!=nil {
	debug("cannot clear tag of the window with id %d: %s\n", id, err)
	return
}
if _, err:=f.Write([]byte(s)); err!=nil {
	debug("cannot write tag of the window with id %d: %s\n", id, err)
	return
}

@
@c
func escapeSymbols(s string) (es string) {
	for _, v:=range s {
		if strings.ContainsRune("|\\/[].+?()*^$", v) {
			es+="\\"
		}
		es+=string(v)
	}
	return
}
