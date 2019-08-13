# The Ri hardware description language
The goal of this project is to develop a hardware description language, exploring some alternative approaches to various problems.
## Ideas
### Basic syntax
```
add() := module {
  a, b : u8 in;
  q : u8 out;
  q = a + b;
};
```
This is a module that's simply an 8-bit adder.
`comb` defines a combinational block, meaning that `q` updates instantly whenever `a` or `b` change.
Module instantiation looks something like this:
```
add_i := add();
add.a = 42;
add.b = 23;
```
### Declarative Metaprogramming
Some form of metaprogramming (code running at compile-time) is a must for any serious HDL.
Imperative (or functional) styles of metaprogramming tend to suffer from poor readability, as the user has to work out the evaluation order (often across multiple modules).
In a declarative style the higher-order form is essentially described as a set of equations and it is up to the compiler to find a solution.
A simple example is type inference, where types propagate through the program.
A more complex example is a module with an indefinite number of ports (for instance, an interconnect).
Other modules can create ports, which then require updates to the internal logic.

Syntax is still TBC.
### Finite State Machine (FSM) Notation
FSMs are the bread and butter of complex hardware designs, as such a special notation is extremely desirable.
For example,
```
fsm {
  q = 0;
  cycle;
  q = 1;
  cycle;
};
```
alternates `q` between 0 and 1.
C-style `if`, `while`, `for`, etc. are supported as well.
### Linking FSMs
Some things are more easily expressed using multiple FSMs in synchrony.
For instance, suppose we want to add flow control to the example above.
```
fsm {
  q_valid = 1;
  q_data = 0;
  while(!q_ready) cycle;
  cycle;
  q_data = 1;
  while(!q_ready) cycle;
  cycle;
};
```
It would be nice if the flow control could be abstracted away, leaving a notation like
```
fsm {
  q.write(0);
  cycle;
  q.write(1);
  cycle;
};
```
In this simple example this could be accomplished by `write` as a "method" like
```
write(x) := {
  valid = 1;
  data = x;
  while(!ready) cycle;
  cycle;
  valid = 0;
};
```
But if we want to do multiple writes in one cycle (to different ports, of course), this will not do.
The linking approach is to define a FSM to describe the handshaking
```
fsm {
  wait sync_write;
  valid = 1;
  while(!ready) cycle;
  cycle;
  valid = 0;
};
```
and a `write` method as
```
write(x) := {
  wait sync_write;
  data = x;
};
```
Here, `sync_write` is a synchronisation point. `wait` ensures that both FSMs need to be interlocked in order to proceed (otherwise, stall cycles are inserted).
With these definitions, writing
```
q.write(0);
r.write(1);
```
will write either to both or to neither.
