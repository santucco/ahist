

/*2:*/


//line ahist.w:8


//line license:1

// This file is part of ahist version 0.1
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



/*6:*/


//line ahist.w:50

"strconv"



/*:6*/



/*14:*/


//line ahist.w:99

"github.com/santucco/goacme"



/*:14*/



/*17:*/


//line ahist.w:123

"strings"



/*:17*/


//line ahist.w:15

)

var(


/*7:*/


//line ahist.w:54

id int



/*:7*/



/*9:*/


//line ahist.w:68

tagname string



/*:9*/



/*18:*/


//line ahist.w:127

name string



/*:18*/



/*39:*/


//line ahist.w:357

histch chan entry= make(chan entry)



/*:39*/


//line ahist.w:19

)

type(


/*38:*/


//line ahist.w:350

entry struct{
b,e int
s string
}



/*:38*/


//line ahist.w:23

)



/*:2*/



/*3:*/


//line ahist.w:27

func main(){


/*10:*/


//line ahist.w:72

tagname= os.Args[0]
if n:=strings.LastIndex(tagname,"/");n!=-1{
tagname= tagname[n:]
}
debug("tagname:%s\n",tagname)



/*:10*/


//line ahist.w:29



/*8:*/


//line ahist.w:58

{
var err error
id,err= strconv.Atoi(os.Getenv("winid"))
if err!=nil{
return
}
}



/*:8*/


//line ahist.w:30



/*16:*/


//line ahist.w:114

w,err:=goacme.Open(id)
if err!=nil{
debug("cannot open a window with id %d: %s\n",id,err)
return
}
defer w.Close()



/*:16*/


//line ahist.w:31



/*11:*/


//line ahist.w:80

{
del:=append([]string{},tagname)
add:=append([]string{},"-"+tagname)
changeTag(w,del,add)
}



/*:11*/


//line ahist.w:32



/*19:*/


//line ahist.w:131

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



/*:19*/


//line ahist.w:33



/*42:*/


//line ahist.w:369

go func(){


/*41:*/


//line ahist.w:365

var hch<-chan*goacme.Event



/*:41*/



/*43:*/


//line ahist.w:383

var h*goacme.Window



/*:43*/



/*47:*/


//line ahist.w:425

var history map[int]int



/*:47*/


//line ahist.w:371

for{
select{
case entr,ok:=<-histch:


/*44:*/


//line ahist.w:387

if!ok{
if h!=nil{
h.Del(true)
h.Close()
h= nil
}
return
}


/*48:*/


//line ahist.w:429

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



/*:48*/


//line ahist.w:396

if ee,ok:=history[entr.b];ok&&ee==entr.e{
continue
}
history[entr.b]= entr.e
debug("writing to the history %d,%d\n",entr.b,entr.e)
h.Write([]byte(fmt.Sprintf("%s:#%d,#%d %q\n",name,entr.b,entr.e,entr.s)))
h.WriteCtl("clean")



/*:44*/


//line ahist.w:375

case ev,ok:=<-hch:


/*45:*/


//line ahist.w:408

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



/*:45*/


//line ahist.w:377

}
}
}()



/*:42*/


//line ahist.w:34



/*15:*/


//line ahist.w:103

for{
ev,err:=w.ReadEvent()
if err!=nil{


/*12:*/


//line ahist.w:88

{
del:=append([]string{},"-"+tagname)
add:=append([]string{},tagname)
changeTag(w,del,add)
}




/*:12*/



/*36:*/


//line ahist.w:339

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:36*/



/*40:*/


//line ahist.w:361

close(histch)



/*:40*/


//line ahist.w:107

return
}


/*20:*/


//line ahist.w:156



/*21:*/


//line ahist.w:164

debug("ev: %#v\n",ev)
s:=""
type_switch:switch{
case ev.Type==goacme.Look|goacme.Tag:


/*22:*/


//line ahist.w:182

s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}


/*28:*/


//line ahist.w:247

if w.WriteCtl("addr=dot")!=nil{


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:249

}
debug("set addr to dot\n")



/*:28*/


//line ahist.w:187




/*:22*/


//line ahist.w:169

case ev.Type==goacme.Look:


/*23:*/


//line ahist.w:192

s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}
b:=ev.Begin
e:=ev.End


/*46:*/


//line ahist.w:421

histch<-entry{b:b,e:e,s:s}



/*:46*/


//line ahist.w:199



/*30:*/


//line ahist.w:270

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:273

}
debug("set addr to %d, %d\n",b,e)



/*:30*/


//line ahist.w:200




/*:23*/


//line ahist.w:171

case ev.Type==goacme.Execute||ev.Type==goacme.Execute|goacme.Tag:


/*24:*/


//line ahist.w:207

switch ev.Text{
case"Look":
s= ev.Arg


/*28:*/


//line ahist.w:247

if w.WriteCtl("addr=dot")!=nil{


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:249

}
debug("set addr to dot\n")



/*:28*/


//line ahist.w:211

break type_switch
case tagname:
continue
case"-"+tagname:


/*12:*/


//line ahist.w:88

{
del:=append([]string{},"-"+tagname)
add:=append([]string{},tagname)
changeTag(w,del,add)
}




/*:12*/



/*36:*/


//line ahist.w:339

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:36*/



/*40:*/


//line ahist.w:361

close(histch)



/*:40*/


//line ahist.w:216

return
}
w.UnreadEvent(ev)
fallthrough




/*:24*/


//line ahist.w:173

case ev.Type==goacme.Insert||ev.Type==goacme.Delete:


/*35:*/


//line ahist.w:320

{
_,_,_,_,d,_,_,_,err:=w.ReadCtl()
if err!=nil{
debug("cannot read from 'ctl' of the window with id %d: %s\n",id,err)
continue
}
debug("dirty: %v\n",d)
var add,del[]string
if d{
add= append(add,"Put")
}else{
del= append(del,"Put")
}
add= append(add,"Undo","Redo")
changeTag(w,del,add)
}



/*:35*/


//line ahist.w:175

continue
default:


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:178

}



/*:21*/


//line ahist.w:157



/*26:*/


//line ahist.w:230

{


/*32:*/


//line ahist.w:294

b,e,err:=w.ReadAddr()
if err!=nil{


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:297

}
debug("read address b: %v, e: %v\n",b,e)



/*:32*/


//line ahist.w:232

if len(s)> 0{


/*31:*/


//line ahist.w:278
{
es:=""
for _,v:=range s{
if strings.ContainsRune("|\\/[].+?()*^$",v){
es+= "\\"
}
es+= string(v)
}
debug("es: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:289

}
}



/*:31*/


//line ahist.w:234

}else{


/*27:*/


//line ahist.w:241
{


/*29:*/


//line ahist.w:254
{
d,err:=w.File("xdata")
if err!=nil{
debug("cannot read from 'xdata' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:258

}

buf:=make([]byte,e-b+1)

for n,_:=d.Read(buf);n> 0;n,_= d.Read(buf){
s+= string(buf[:n])
}
debug("read address from xdata b: %v, e: %v\n",b,e)
}



/*:29*/


//line ahist.w:242



/*31:*/


//line ahist.w:278
{
es:=""
for _,v:=range s{
if strings.ContainsRune("|\\/[].+?()*^$",v){
es+= "\\"
}
es+= string(v)
}
debug("es: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:289

}
}



/*:31*/


//line ahist.w:243

}



/*:27*/


//line ahist.w:236

}
}



/*:26*/


//line ahist.w:158



/*32:*/


//line ahist.w:294

b,e,err:=w.ReadAddr()
if err!=nil{


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:297

}
debug("read address b: %v, e: %v\n",b,e)



/*:32*/


//line ahist.w:159



/*34:*/


//line ahist.w:310



/*33:*/


//line ahist.w:302

if w.WriteCtl("dot=addr\nshow")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:305

}
debug("set dot to addr\n")



/*:33*/


//line ahist.w:311

if w.WriteCtl("show")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*25:*/


//line ahist.w:224

w.UnreadEvent(ev)
continue



/*:25*/


//line ahist.w:314

}
debug("show dot\n")



/*:34*/


//line ahist.w:160



/*46:*/


//line ahist.w:421

histch<-entry{b:b,e:e,s:s}



/*:46*/


//line ahist.w:161




/*:20*/


//line ahist.w:110

}



/*:15*/


//line ahist.w:35

}



/*:3*/



/*5:*/


//line ahist.w:44

func debug(f string,args...interface{}){
//	fmt.Fprintf(os.Stderr, f, args...)
}



/*:5*/



/*49:*/


//line ahist.w:446

func changeTag(w*goacme.Window,del[]string,add[]string){
if add==nil&&del==nil{
return
}


/*50:*/


//line ahist.w:458

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



/*:50*/


//line ahist.w:451



/*51:*/


//line ahist.w:477

if n= strings.LastIndex(s,"|");n==-1{
n= 0
}else{
n++
}
s= s[n:]
s= strings.TrimLeft(s," ")
tag:=strings.Split(s," ")



/*:51*/


//line ahist.w:452



/*52:*/


//line ahist.w:488

newtag:=append([]string{},"")


/*53:*/


//line ahist.w:496

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



/*:53*/


//line ahist.w:490



/*54:*/


//line ahist.w:509

for _,v:=range add{
for i:=0;i<len(tag);{
if tag[i]!=v{
i++
continue
}
copy(tag[i:],tag[i+1:])
tag= tag[:len(tag)-1]
}
}



/*:54*/


//line ahist.w:491

newtag= append(newtag,add...)
newtag= append(newtag,tag...)



/*:52*/


//line ahist.w:453



/*55:*/


//line ahist.w:522

s= strings.Join(newtag," ")
if err:=w.WriteCtl("cleartag");err!=nil{
debug("cannot clear tag of the window with id %d: %s\n",id,err)
return
}
if _,err:=f.Write([]byte(s));err!=nil{
debug("cannot write tag of the window with id %d: %s\n",id,err)
return
}

/*:55*/


//line ahist.w:454

}



/*:49*/


