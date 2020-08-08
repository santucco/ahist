

/*2:*/


//line ahist.w:8


//line license:1

// This file is part of ahist
//
// Copyright (c) 2020 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//line ahist.w:11

package main

import(


/*4:*/


//line ahist.w:39

"fmt"
"os"



/*:4*/



/*9:*/


//line ahist.w:66

"strconv"



/*:9*/



/*17:*/


//line ahist.w:115

"github.com/santucco/goacme"



/*:17*/



/*20:*/


//line ahist.w:139

"strings"



/*:20*/


//line ahist.w:15

)

var(


/*5:*/


//line ahist.w:44

dbg bool



/*:5*/



/*10:*/


//line ahist.w:70

id int



/*:10*/



/*12:*/


//line ahist.w:84

tagname string



/*:12*/



/*21:*/


//line ahist.w:143

name string



/*:21*/



/*34:*/


//line ahist.w:301

lentr entry



/*:34*/



/*44:*/


//line ahist.w:405

histch chan entry= make(chan entry)



/*:44*/


//line ahist.w:19

)

type(


/*43:*/


//line ahist.w:398

entry struct{
b,e int
s string
}



/*:43*/


//line ahist.w:23

)



/*:2*/



/*3:*/


//line ahist.w:27

func main(){


/*13:*/


//line ahist.w:88

tagname= os.Args[0]
if n:=strings.LastIndex(tagname,"/");n!=-1{
tagname= tagname[n:]
}
debug("tagname:%s\n",tagname)



/*:13*/


//line ahist.w:29



/*11:*/


//line ahist.w:74

{
var err error
id,err= strconv.Atoi(os.Getenv("winid"))
if err!=nil{
return
}
}



/*:11*/


//line ahist.w:30



/*19:*/


//line ahist.w:130

w,err:=goacme.Open(id)
if err!=nil{
debug("cannot open a window with id %d: %s\n",id,err)
return
}
defer w.Close()



/*:19*/


//line ahist.w:31



/*14:*/


//line ahist.w:96

{
del:=[]string{tagname,"-"+tagname,"-"+tagname+"+","-"+tagname+"-"}
add:=[]string{"-"+tagname}
changeTag(w,del,add)
}



/*:14*/


//line ahist.w:32



/*22:*/


//line ahist.w:147

{
f,err:=w.File("tag")
if err!=nil{
debug("cannot read from 'tag' of the window with id %d: %s\n",id,err)
return
}
if _,err:=f.Seek(0,0);err!=nil{
debug("cannot seek to the start 'tag' of the window with id %d: %s\n",id,err)
return
}
var b[1000]byte
n,err:=f.Read(b[:])
if err!=nil{
debug("cannot read tag of the window with id %d: %s\n",id,err)
return
}
ss:=strings.Split(string(b[:n])," ")
if len(ss)==0{
return
}
name= string(ss[0])
}



/*:22*/


//line ahist.w:33



/*47:*/


//line ahist.w:417

go func(){


/*46:*/


//line ahist.w:413

var hch<-chan*goacme.Event



/*:46*/



/*48:*/


//line ahist.w:431

var h*goacme.Window



/*:48*/



/*52:*/


//line ahist.w:474

var history map[int]int



/*:52*/


//line ahist.w:419

for{
select{
case entr,ok:=<-histch:


/*49:*/


//line ahist.w:435

if!ok{
if h!=nil{
h.Del(true)
h.Close()
h= nil
}
return
}


/*53:*/


//line ahist.w:478

if h==nil{
var err error
if h,err= goacme.New();err!=nil{
return
}
h.WriteCtl("name %s",name+"+History")
if hch,err= h.EventChannel(1,goacme.AllTypes);err!=nil{
return
}
history= make(map[int]int)
}



/*:53*/


//line ahist.w:444

if ee,ok:=history[entr.b];ok&&ee==entr.e{
continue
}
history[entr.b]= entr.e
debug("writing to the history %d,%d\n",entr.b,entr.e)
h.Write([]byte(fmt.Sprintf("%s:#%d,#%d %q\n",name,entr.b,entr.e,entr.s)))
h.WriteCtl("clean")



/*:49*/


//line ahist.w:423

case ev,ok:=<-hch:


/*50:*/


//line ahist.w:456

if!ok{
debug("history is closed\n")
h.Del(true)
h.Close()
h= nil
hch= nil
history= nil
continue
}
h.UnreadEvent(ev)



/*:50*/


//line ahist.w:425

}
}
}()



/*:47*/


//line ahist.w:34



/*18:*/


//line ahist.w:119

for{
ev,err:=w.ReadEvent()
if err!=nil{


/*15:*/


//line ahist.w:104

{
del:=[]string{tagname,"-"+tagname,"-"+tagname+"+","-"+tagname+"-"}
add:=[]string{tagname}
changeTag(w,del,add)
}




/*:15*/



/*41:*/


//line ahist.w:387

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:41*/



/*45:*/


//line ahist.w:409

close(histch)



/*:45*/


//line ahist.w:123

return
}


/*23:*/


//line ahist.w:172



/*24:*/


//line ahist.w:180

debug("incoming event: %+v\n",ev)
s:=""
type_switch:switch{
case ev.Type==goacme.Look|goacme.Tag:


/*25:*/


//line ahist.w:198

s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}


/*31:*/


//line ahist.w:268

if w.WriteCtl("addr=dot")!=nil{


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:270

}
debug("set addr to dot\n")



/*:31*/


//line ahist.w:203




/*:25*/


//line ahist.w:185

case ev.Type==goacme.Look:


/*26:*/


//line ahist.w:206

s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}
b:=ev.Begin
e:=ev.End


/*33:*/


//line ahist.w:291

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:294

}
debug("set addr to %d, %d\n",b,e)



/*:33*/


//line ahist.w:213




/*:26*/


//line ahist.w:187

case ev.Type==goacme.Execute||ev.Type==goacme.Execute|goacme.Tag:


/*27:*/


//line ahist.w:221

switch strings.TrimSpace(ev.Text){
case"Look":
s= ev.Arg


/*31:*/


//line ahist.w:268

if w.WriteCtl("addr=dot")!=nil{


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:270

}
debug("set addr to dot\n")



/*:31*/


//line ahist.w:225

break type_switch
case tagname:
continue
case"-"+tagname:
debug("exiting\n")


/*15:*/


//line ahist.w:104

{
del:=[]string{tagname,"-"+tagname,"-"+tagname+"+","-"+tagname+"-"}
add:=[]string{tagname}
changeTag(w,del,add)
}




/*:15*/



/*41:*/


//line ahist.w:387

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:41*/



/*45:*/


//line ahist.w:409

close(histch)



/*:45*/


//line ahist.w:231

return
case tagname+"+":


/*6:*/


//line ahist.w:48

dbg= true
debug("debug has been switched on\n")



/*:6*/


//line ahist.w:234

continue
case tagname+"-":


/*7:*/


//line ahist.w:53

debug("debug has been switched off\n")
dbg= false



/*:7*/


//line ahist.w:237

continue
}
w.UnreadEvent(ev)
fallthrough




/*:27*/


//line ahist.w:189

case ev.Type==goacme.Insert||ev.Type==goacme.Delete:


/*40:*/


//line ahist.w:369

{
_,_,_,_,d,_,_,_,err:=w.ReadCtl()
if err!=nil{
debug("cannot read from 'ctl' of the window with id %d: %s\n",id,err)
continue
}
debug("dirty: %v\n",d)
del:=[]string{"Put","Undo","Redo"}
var add[]string
if d{
add= append(add,"Put")
}
add= append(add,"Undo","Redo")
changeTag(w,del,add)
}



/*:40*/


//line ahist.w:191

continue
default:


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:194

}



/*:24*/


//line ahist.w:173



/*29:*/


//line ahist.w:251

{


/*37:*/


//line ahist.w:343

b,e,err:=w.ReadAddr()
if err!=nil{


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:346

}
debug("read address b: %v, e: %v\n",b,e)



/*:37*/


//line ahist.w:253

if len(s)> 0{


/*36:*/


//line ahist.w:315
{
debug("last entry : %v\n",lentr)
if len(s)==0{
if!lentr.empty(){
b= lentr.b
e= lentr.e
s= lentr.s


/*33:*/


//line ahist.w:291

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:294

}
debug("set addr to %d, %d\n",b,e)



/*:33*/


//line ahist.w:322

}
}else if b!=e{
lentr= entry{b,e,s}


/*51:*/


//line ahist.w:469

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:51*/


//line ahist.w:326

}
es:=""
for _,v:=range s{
if strings.ContainsRune("|\\/[].+?()*^$",v){
es+= "\\"
}
es+= string(v)
}
debug("escaped search string: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:338

}
}



/*:36*/


//line ahist.w:255

}else{


/*30:*/


//line ahist.w:262
{


/*32:*/


//line ahist.w:275
{
d,err:=w.File("xdata")
if err!=nil{
debug("cannot read from 'xdata' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:279

}

buf:=make([]byte,e-b+1)

for n,_:=d.Read(buf);n> 0;n,_= d.Read(buf){
s+= string(buf[:n])
}
debug("read address from xdata b: %v, e: %v\n",b,e)
}



/*:32*/


//line ahist.w:263



/*36:*/


//line ahist.w:315
{
debug("last entry : %v\n",lentr)
if len(s)==0{
if!lentr.empty(){
b= lentr.b
e= lentr.e
s= lentr.s


/*33:*/


//line ahist.w:291

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:294

}
debug("set addr to %d, %d\n",b,e)



/*:33*/


//line ahist.w:322

}
}else if b!=e{
lentr= entry{b,e,s}


/*51:*/


//line ahist.w:469

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:51*/


//line ahist.w:326

}
es:=""
for _,v:=range s{
if strings.ContainsRune("|\\/[].+?()*^$",v){
es+= "\\"
}
es+= string(v)
}
debug("escaped search string: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:338

}
}



/*:36*/


//line ahist.w:264

}



/*:30*/


//line ahist.w:257

}
}



/*:29*/


//line ahist.w:174



/*37:*/


//line ahist.w:343

b,e,err:=w.ReadAddr()
if err!=nil{


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:346

}
debug("read address b: %v, e: %v\n",b,e)



/*:37*/


//line ahist.w:175



/*39:*/


//line ahist.w:359



/*38:*/


//line ahist.w:351

if w.WriteCtl("dot=addr\nshow")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:354

}
debug("set dot to addr\n")



/*:38*/


//line ahist.w:360

if w.WriteCtl("show")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*28:*/


//line ahist.w:245

w.UnreadEvent(ev)
continue



/*:28*/


//line ahist.w:363

}
debug("show dot\n")



/*:39*/


//line ahist.w:176



/*51:*/


//line ahist.w:469

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:51*/


//line ahist.w:177




/*:23*/


//line ahist.w:126

}



/*:18*/


//line ahist.w:35

}



/*:3*/



/*8:*/


//line ahist.w:58

func debug(f string,args...interface{}){
if dbg{
fmt.Fprintf(os.Stderr,f,args...)
}
}



/*:8*/



/*35:*/


//line ahist.w:305

func(this entry)empty()bool{
return this.b==this.e
}



/*:35*/



/*54:*/


//line ahist.w:495

func changeTag(w*goacme.Window,del[]string,add[]string){
if add==nil&&del==nil{
return
}


/*55:*/


//line ahist.w:507

f,err:=w.File("tag")
if err!=nil{
debug("cannot read from 'tag' of the window with id %d: %s\n",id,err)
return
}
if _,err:=f.Seek(0,0);err!=nil{
debug("cannot seek to the start 'tag' of the window with id %d: %s\n",id,err)
return
}
var b[1000]byte
n,err:=f.Read(b[:])
if err!=nil{
debug("cannot read tag of the window with id %d: %s\n",id,err)
return
}
s:=string(b[:n])



/*:55*/


//line ahist.w:500



/*56:*/


//line ahist.w:526

if n= strings.LastIndex(s,"|");n==-1{
n= 0
}else{
n++
}
s= s[n:]
s= strings.TrimLeft(s," ")
tag:=strings.Split(s," ")



/*:56*/


//line ahist.w:501



/*57:*/


//line ahist.w:537

newtag:=append([]string{},"")


/*58:*/


//line ahist.w:544

for _,v:=range del{
for i:=0;i<len(tag);{
if tag[i]!=v{
i++
continue
}
copy(tag[i:],tag[i+1:])
tag= tag[:len(tag)-1]
}
}



/*:58*/


//line ahist.w:539

newtag= append(newtag,add...)
newtag= append(newtag,tag...)



/*:57*/


//line ahist.w:502



/*59:*/


//line ahist.w:557

s= strings.Join(newtag," ")
if err:=w.WriteCtl("cleartag");err!=nil{
debug("cannot clear tag of the window with id %d: %s\n",id,err)
return
}
if _,err:=f.Write([]byte(s));err!=nil{
debug("cannot write tag of the window with id %d: %s\n",id,err)
return
}

/*:59*/


//line ahist.w:503

}



/*:54*/


