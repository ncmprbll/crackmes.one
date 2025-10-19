# _starting off

_starting off research

# TLDR

This binary can be exploited in the following way:

```
$ ./chall {1..65572}
Entry, main
flag{done}
Segmentation fault (core dumped)
```

# Long version

Upon first examination, we have what appears to be regular binary. However, it crashes whenever we run it:

```
$ ./chall
Entry, main
Segmentation fault (core dumped)
```

Passing arguments yields the same result:

```
$ ./chall 1
Entry, main
Segmentation fault (core dumped)
$ ./chall 2
Entry, main
Segmentation fault (core dumped)
$ ./chall 1 2
Entry, main
Segmentation fault (core dumped)
$ ./chall flag.txt
Entry, main
Segmentation fault (core dumped)
```

`Entry, main` gives us no insight whatsoever.

Let's hop on `gdb` and see where it takes us. We'll start by using `starti` command to start the program and immediately break after executing first instruction:

```
(gdb) starti
Starting program: .../68eb89c32d267f28f69b7544/todo/chall

This GDB supports auto-downloading debuginfo from the following URLs:
  <https://debuginfod.ubuntu.com>
Enable debuginfod for this session? (y or [n]) n
Debuginfod has been disabled.
To make this setting permanent, add 'set debuginfod enabled off' to .gdbinit.

Program stopped.
0x000000000001007a in ?? ()
(gdb)
```

`0x000000000001007a` appears to be program's `_start` - an entry point. Before going any further, we'll also examine executable's address space using `info proc mappings`:

```
(gdb) info proc mappings
process 1565
Mapped address spaces:

          Start Addr           End Addr       Size     Offset  Perms  objfile
             0x10000            0x11000     0x1000     0x1000  r-xp   .../68eb89c32d267f28f69b7544/todo/chall
      0x7ffff7ff9000     0x7ffff7ffd000     0x4000        0x0  r--p   [vvar]
      0x7ffff7ffd000     0x7ffff7fff000     0x2000        0x0  r-xp   [vdso]
      0x7ffffffde000     0x7ffffffff000    0x21000        0x0  rw-p   [stack]
(gdb)
```

This output tells us that the program's address space is between `0x10000` and `0x11000` and we will keep it in mind for now. We are not interested in `vvar`, `vdso` and `stack` shared libraries probably exported by the kernel.

Let's continue our journey by using `x/24i $pc` command and find the reason for our segmentation fault:

```
(gdb) x/24i $pc
=> 0x1007a:     push   rbp
   0x1007b:     mov    rbp,rsp
   0x1007e:     sub    rsp,0x30
   0x10082:     movabs rax,0x6d202c7972746e45
   0x1008c:     mov    edx,0xa6e6961
   0x10091:     mov    QWORD PTR [rbp-0x30],rax
   0x10095:     mov    QWORD PTR [rbp-0x28],rdx
   0x10099:     mov    QWORD PTR [rbp-0x20],0x0
   0x100a1:     mov    QWORD PTR [rbp-0x18],0x0
   0x100a9:     mov    DWORD PTR [rbp-0x4],0xc
   0x100b0:     mov    edx,DWORD PTR [rbp-0x4]
   0x100b3:     lea    rax,[rbp-0x30]
   0x100b7:     mov    esi,edx
   0x100b9:     mov    rdi,rax
   0x100bc:     call   0x10000
   0x100c1:     mov    eax,0x0
   0x100c6:     leave
   0x100c7:     ret
   0x100c8:     data16 ins BYTE PTR es:[rdi],dx
   0x100ca:     (bad)
   0x100cb:     addr32 cs je 0x10147
   0x100cf:     je     0x100d1
   0x100d1:     add    BYTE PTR [rax],al
   0x100d3:     add    BYTE PTR [rax],al
(gdb)
```

Nothing seems out of order. Let's go to `0x100bc` and look inside of the function call:

```
(gdb) b *0x100bc
Breakpoint 1 at 0x100bc
(gdb) c
Continuing.

Breakpoint 1, 0x00000000000100bc in ?? ()
(gdb) ni
Entry, main
0x00000000000100c1 in ?? ()
(gdb) x/12i $pc
=> 0x100c1:     mov    eax,0x0
   0x100c6:     leave
   0x100c7:     ret
   0x100c8:     data16 ins BYTE PTR es:[rdi],dx
   0x100ca:     (bad)
   0x100cb:     addr32 cs je 0x10147
   0x100cf:     je     0x100d1
   0x100d1:     add    BYTE PTR [rax],al
   0x100d3:     add    BYTE PTR [rax],al
   0x100d5:     add    BYTE PTR [rax],al
   0x100d7:     add    BYTE PTR [rax+rax*1],dl
   0x100da:     add    BYTE PTR [rax],al
(gdb) ni
0x00000000000100c6 in ?? ()
(gdb) ni
0x00000000000100c7 in ?? ()
(gdb) ni
0x0000000000000001 in ?? ()
(gdb) ni

Program received signal SIGSEGV, Segmentation fault.
0x0000000000000001 in ?? ()
(gdb)
```

We crash after executing an instruction at a return address of a `ret` instruction. It looks like `leave` is responsible for overriding our return address. If we run the program again and follow it to the `leave` instruction, we can run `info reg rsp rbp` and `x/g $rsp` commands to examine our stack pointer, frame pointer and our overriden return address:

```
(gdb) b *0x100c6
Breakpoint 4 at 0x100c6
(gdb) c
Continuing.

Breakpoint 1, 0x00000000000100bc in ?? ()
(gdb) c
Continuing.
Entry, main

Breakpoint 4, 0x00000000000100c6 in ?? ()
(gdb) x/2i $pc
=> 0x100c6:     leave
   0x100c7:     ret
(gdb) info reg rsp rbp
rsp            0x7fffffffe0b8      0x7fffffffe0b8
rbp            0x7fffffffe0e8      0x7fffffffe0e8
(gdb) ni
0x00000000000100c7 in ?? ()
(gdb) info reg rsp rbp
rsp            0x7fffffffe0f0      0x7fffffffe0f0
rbp            0x0                 0x0
(gdb) x/g $rsp
0x7fffffffe0f0: 0x0000000000000001
(gdb) ni
0x0000000000000001 in ?? ()
(gdb) ni

Program received signal SIGSEGV, Segmentation fault.
0x0000000000000001 in ?? ()
(gdb)
```

`0x0000000000000001` is the culprit. But what is it? If we run `run` command with arguments, we can actually make an acute observation:

```
(gdb) run "1" "2" "3"
The program being debugged has been started already.
Start it from the beginning? (y or n) y
Starting program: .../68eb89c32d267f28f69b7544/todo/chall "1" "2" "3"

Breakpoint 1, 0x00000000000100bc in ?? ()
(gdb) c
Continuing.
Entry, main

Breakpoint 4, 0x00000000000100c6 in ?? ()
(gdb) c
Continuing.

Program received signal SIGSEGV, Segmentation fault.
0x0000000000000004 in ?? ()
(gdb)
```

It's `0x0000000000000004` now. Our return address is the amount of arguments or `argc`! Three explicit arguments "1", "2" and "3" and one implict argument - path to the executable. Now we just have to find a function to point at. Using commands `find` and `x` we can find if `flag.txt` string exists in the file, and if it does, then we can find references from the code:

```
(gdb) find 0x10000, +0x1000, "flag.txt"
0x100c8
1 pattern found.
(gdb) find 0x10000, +0x1000, 0x100c8
0x10030
0x10280
2 patterns found.
(gdb) x/32i 0x10030-0xc
   0x10024:     push   rbp
   0x10025:     mov    rbp,rsp
   0x10028:     sub    rsp,0x50
   0x1002c:     mov    QWORD PTR [rbp-0x8],0x100c8
   0x10034:     mov    rdi,QWORD PTR [rbp-0x8]
   0x10038:     mov    rax,0x2
   0x1003f:     mov    rsi,0x2
   0x10046:     mov    rdx,0x1ff
   0x1004d:     syscall
   0x1004f:     mov    rdi,rax
   0x10052:     lea    rsi,[rbp-0x50]
   0x10056:     mov    rax,0x0
   0x1005d:     mov    rdx,0xb
   0x10064:     syscall
   0x10066:     lea    rax,[rbp-0x50]
   0x1006a:     mov    esi,0xb
   0x1006f:     mov    rdi,rax
   0x10072:     call   0x10000
   0x10077:     nop
   0x10078:     leave
   0x10079:     ret
   0x1007a:     push   rbp
   0x1007b:     mov    rbp,rsp
   0x1007e:     sub    rsp,0x30
   0x10082:     movabs rax,0x6d202c7972746e45
   0x1008c:     mov    edx,0xa6e6961
   0x10091:     mov    QWORD PTR [rbp-0x30],rax
   0x10095:     mov    QWORD PTR [rbp-0x28],rdx
   0x10099:     mov    QWORD PTR [rbp-0x20],0x0
   0x100a1:     mov    QWORD PTR [rbp-0x18],0x0
   0x100a9:     mov    DWORD PTR [rbp-0x4],0xc
   0x100b0:     mov    edx,DWORD PTR [rbp-0x4]
(gdb)
```

We found function prologue at `0x10024` which references `flag.txt` string at `0x1002c` (`mov QWORD PTR [rbp-0x8], 0x100c8`). This function also makes 2 syscalls with numbers `0x0` (read) and `0x2` (open). But how do we exploit the binary to point to this function? Do we... supply `0x10024` (`65572` in decimal) arguments? We could, of course. But we could also use bash's brace expansion:

```
$ ./chall {1..65572}
Entry, main
flag{done}
Segmentation fault (core dumped)
```