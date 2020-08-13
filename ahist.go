

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


//line ahist.w:140

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


//line ahist.w:144

name string



/*:21*/



/*35:*/


//line ahist.w:313

lentr entry



/*:35*/



/*45:*/


//line ahist.w:411

histch chan entry= make(chan entry)



/*:45*/


//line ahist.w:19

)

type(


/*44:*/


//line ahist.w:404

entry struct{
b,e int
s string
}



/*:44*/


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


//line ahist.w:131

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


//line ahist.w:148

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



/*48:*/


//line ahist.w:423

go func(){


/*47:*/


//line ahist.w:419

var hch<-chan*goacme.Event



/*:47*/



/*49:*/


//line ahist.w:437

var h*goacme.Window



/*:49*/



/*54:*/


//line ahist.w:496

var history map[int]int



/*:54*/


//line ahist.w:425

for{
select{
case entr,ok:=<-histch:


/*50:*/


//line ahist.w:441

if!ok{
if h!=nil{
h.Del(true)
h.Close()
h= nil
}
return
}


/*55:*/


//line ahist.w:500

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



/*:55*/


//line ahist.w:450

if history[entr.b]!=entr.e{
history[entr.b]= entr.e
debug("writing to the history %d,%d\n",entr.b,entr.e)
h.Write([]byte(fmt.Sprintf("%s:#%d,#%d %q\n",name,entr.b,entr.e,entr.s)))
h.WriteCtl("clean")
}
debug("selecting the current position #%d,#%d in the history\n",entr.b,entr.e)
es:=fmt.Sprintf("#%d,#%d",entr.b,entr.e)


/*52:*/


//line ahist.w:483

if err:=h.WriteAddr("/%s/-+",es);err!=nil{
debug("writing of addr failed: %s\n",err)
}else if err:=h.WriteCtl("dot=addr\nshow");err!=nil{
debug("writing of ctl failed: %s\n",err)
}



/*:52*/


//line ahist.w:459




/*:50*/


//line ahist.w:429

case ev,ok:=<-hch:


/*51:*/


//line ahist.w:464

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
if ev.Type==goacme.Look{
debug("incoming event: %+v\n",ev)
debug("selecting the current position %q in the history\n",ev.Text)
es:=escapeSymbols(ev.Text)


/*52:*/


//line ahist.w:483

if err:=h.WriteAddr("/%s/-+",es);err!=nil{
debug("writing of addr failed: %s\n",err)
}else if err:=h.WriteCtl("dot=addr\nshow");err!=nil{
debug("writing of ctl failed: %s\n",err)
}



/*:52*/


//line ahist.w:479

}



/*:51*/


//line ahist.w:431

}
}
}()



/*:48*/


//line ahist.w:34



/*18:*/


//line ahist.w:119



/*41:*/


//line ahist.w:375

{
_,_,_,_,d,_,_,_,err:=w.ReadCtl()
if err!=nil{
debug("cannot read from 'ctl' of the window with id %d: %s\n",id,err)
}else{
debug("dirty: %v\n",d)
del:=[]string{"Put","Undo","Redo"}
var add[]string
if d{
add= append(add,"Put")
}
add= append(add,"Undo","Redo")
changeTag(w,del,add)
}
}



/*:41*/


//line ahist.w:120

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



/*42:*/


//line ahist.w:393

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:42*/



/*46:*/


//line ahist.w:415

close(histch)



/*:46*/


//line ahist.w:124

return
}


/*23:*/


//line ahist.w:173



/*24:*/


//line ahist.w:181

debug("incoming event: %+v\n",ev)
s:=""
b:=ev.Begin
e:=ev.End
type_switch:switch{
case ev.Type==goacme.Look|goacme.Tag:


/*25:*/


//line ahist.w:202

b,e= 0,0
s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}


/*32:*/


//line ahist.w:280

if w.WriteCtl("addr=dot")!=nil{


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:282

}
debug("set addr to dot\n")



/*:32*/


//line ahist.w:208




/*:25*/


//line ahist.w:188

case ev.Type==goacme.Look:


/*26:*/


//line ahist.w:211

s= ev.Text
if len(ev.Arg)> 0{
s+= " "+ev.Arg
}


/*34:*/


//line ahist.w:303

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:306

}
debug("set addr to %d, %d\n",b,e)



/*:34*/


//line ahist.w:216




/*:26*/


//line ahist.w:190

case ev.Type==goacme.Execute||ev.Type==goacme.Execute|goacme.Tag:


/*27:*/


//line ahist.w:224

switch strings.TrimSpace(ev.Text){
case"Look":


/*28:*/


//line ahist.w:251

s= ev.Arg


/*32:*/


//line ahist.w:280

if w.WriteCtl("addr=dot")!=nil{


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:282

}
debug("set addr to dot\n")



/*:32*/


//line ahist.w:253



/*38:*/


//line ahist.w:349

b,e,err= w.ReadAddr()
if err!=nil{


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:352

}
debug("read address b: %v, e: %v\n",b,e)



/*:38*/


//line ahist.w:254




/*:28*/


//line ahist.w:227

break type_switch
case tagname:
continue
case"-"+tagname+"+":
fallthrough
case"-"+tagname+"-":
fallthrough
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



/*42:*/


//line ahist.w:393

{
del:=append([]string{},"Put","Undo","Redo")
changeTag(w,del,nil)
}



/*:42*/



/*46:*/


//line ahist.w:415

close(histch)



/*:46*/


//line ahist.w:237

return
case tagname+"+":


/*6:*/


//line ahist.w:48

dbg= true
debug("debug has been switched on\n")



/*:6*/


//line ahist.w:240

continue
case tagname+"-":


/*7:*/


//line ahist.w:53

debug("debug has been switched off\n")
dbg= false



/*:7*/


//line ahist.w:243

continue
}
w.UnreadEvent(ev)
fallthrough



/*:27*/


//line ahist.w:192

case ev.Type==goacme.Insert||ev.Type==goacme.Delete:


/*41:*/


//line ahist.w:375

{
_,_,_,_,d,_,_,_,err:=w.ReadCtl()
if err!=nil{
debug("cannot read from 'ctl' of the window with id %d: %s\n",id,err)
}else{
debug("dirty: %v\n",d)
del:=[]string{"Put","Undo","Redo"}
var add[]string
if d{
add= append(add,"Put")
}
add= append(add,"Undo","Redo")
changeTag(w,del,add)
}
}



/*:41*/


//line ahist.w:194

continue
default:


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:197

}



/*:24*/


//line ahist.w:174



/*30:*/


//line ahist.w:263

{
if len(s)> 0{


/*37:*/


//line ahist.w:327
{
debug("last entry : %v\n",lentr)
if len(s)==0{
if!lentr.empty(){
b= lentr.b
e= lentr.e
s= lentr.s


/*34:*/


//line ahist.w:303

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:306

}
debug("set addr to %d, %d\n",b,e)



/*:34*/


//line ahist.w:334

}
}else if b!=e{
lentr= entry{b,e,s}


/*53:*/


//line ahist.w:491

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:53*/


//line ahist.w:338

}
es:=escapeSymbols(s)
debug("escaped search string: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:344

}
}



/*:37*/


//line ahist.w:266

}else{


/*31:*/


//line ahist.w:274
{


/*33:*/


//line ahist.w:287
{
d,err:=w.File("xdata")
if err!=nil{
debug("cannot read from 'xdata' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:291

}

buf:=make([]byte,e-b+1)

for n,_:=d.Read(buf);n> 0;n,_= d.Read(buf){
s+= string(buf[:n])
}
debug("read address from xdata b: %v, e: %v\n",b,e)
}



/*:33*/


//line ahist.w:275



/*37:*/


//line ahist.w:327
{
debug("last entry : %v\n",lentr)
if len(s)==0{
if!lentr.empty(){
b= lentr.b
e= lentr.e
s= lentr.s


/*34:*/


//line ahist.w:303

if err:=w.WriteAddr("#%d,#%d",b,e);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:306

}
debug("set addr to %d, %d\n",b,e)



/*:34*/


//line ahist.w:334

}
}else if b!=e{
lentr= entry{b,e,s}


/*53:*/


//line ahist.w:491

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:53*/


//line ahist.w:338

}
es:=escapeSymbols(s)
debug("escaped search string: %q\n",es)
if err:=w.WriteAddr("/%s/",es);err!=nil{
debug("cannot write to 'addr' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:344

}
}



/*:37*/


//line ahist.w:276

}



/*:31*/


//line ahist.w:268



/*38:*/


//line ahist.w:349

b,e,err= w.ReadAddr()
if err!=nil{


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:352

}
debug("read address b: %v, e: %v\n",b,e)



/*:38*/


//line ahist.w:269

}
}



/*:30*/


//line ahist.w:175



/*38:*/


//line ahist.w:349

b,e,err= w.ReadAddr()
if err!=nil{


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:352

}
debug("read address b: %v, e: %v\n",b,e)



/*:38*/


//line ahist.w:176



/*40:*/


//line ahist.w:365



/*39:*/


//line ahist.w:357

if w.WriteCtl("dot=addr")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:360

}
debug("set dot to addr\n")



/*:39*/


//line ahist.w:366

if w.WriteCtl("show")!=nil{
debug("cannot write to 'ctl' of the window with id %d: %s\n",id,err)


/*29:*/


//line ahist.w:257

w.UnreadEvent(ev)
continue



/*:29*/


//line ahist.w:369

}
debug("show dot\n")



/*:40*/


//line ahist.w:177



/*53:*/


//line ahist.w:491

debug("request to store a history: %v,%v %q\n",b,e,s)
histch<-entry{b:b,e:e,s:s}



/*:53*/


//line ahist.w:178




/*:23*/


//line ahist.w:127

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



/*36:*/


//line ahist.w:317

func(this entry)empty()bool{
return this.b==this.e
}



/*:36*/



/*56:*/


//line ahist.w:517

func changeTag(w*goacme.Window,del[]string,add[]string){
if add==nil&&del==nil{
return
}


/*57:*/


//line ahist.w:529

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



/*:57*/


//line ahist.w:522



/*58:*/


//line ahist.w:548

if n= strings.LastIndex(s,"|");n==-1{
n= 0
}else{
n++
}
s= s[n:]
s= strings.TrimLeft(s," ")
tag:=strings.Split(s," ")



/*:58*/


//line ahist.w:523



/*59:*/


//line ahist.w:559

newtag:=append([]string{},"")


/*60:*/


//line ahist.w:566

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



/*:60*/


//line ahist.w:561

newtag= append(newtag,add...)
newtag= append(newtag,tag...)



/*:59*/


//line ahist.w:524



/*61:*/


//line ahist.w:579

s= strings.Join(newtag," ")
if err:=w.WriteCtl("cleartag");err!=nil{
debug("cannot clear tag of the window with id %d: %s\n",id,err)
return
}
if _,err:=f.Write([]byte(s));err!=nil{
debug("cannot write tag of the window with id %d: %s\n",id,err)
return
}



/*:61*/


//line ahist.w:525

}



/*:56*/



/*62:*/


//line ahist.w:591

func escapeSymbols(s string)(es string){
for _,v:=range s{
if strings.ContainsRune("|\\/[].+?()*^$",v){
es+= "\\"
}
es+= string(v)
}
return
}

/*:62*/


