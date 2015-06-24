Goop
====

Description
-----------

The Goop (Go Object-Oriented Programming) package provides support for dynamic object-oriented programming constructs in Go, much like those that appear in various scripting languages.  The goal is to integrate fast, native-Go objects and slower but more flexible Goop objects within the same program.

Features
--------

Goop provides the following features, which are borrowed from an assortment of object-oriented programming languages:

* a [prototype-based object model](http://en.wikipedia.org/wiki/Prototype-based_programming)

* support for both *ex nihilo* and constructor-based object creation

* the ability to add, replace, and remove data fields and method functions at will

* multiple inheritance

* dynamically modifiable inheritance hierarchies (even on a per-object basis)

* type-dependent dispatch (i.e., multiple methods with the same name but different argument types)

Installation
------------

Install Goop with [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies):

<pre>
    go get github.com/losalamos/goop
</pre>

Documentation
-------------

Pre-built documentation for the Goop API is available online at <http://godoc.org/github.com/losalamos/goop>, courtesy of [GoDoc](http://godoc.org/).

Once you install Goop, you can view the API locally with [`godoc`](http://golang.org/cmd/godoc/), for example by running

<pre>
    godoc github.com/losalamos/goop
</pre>

to display the Goop documentation on screen or by running

<pre>
    godoc -http=:6060
</pre>

to start a local Web server then viewing the HTML-formatted documentation at <http://localhost:6060/pkg/github.com/losalamos/goop/> in your favorite browser.

Performance
-----------

Goop is unfortunately extremely slow.  Goop programs have to pay for
their flexibility in terms of performance.  To determine just how bad
the performance is on your computer, you can run the microbenchmarks
included in
[`goop_test.go`](http://github.com/losalamos/goop/blob/master/goop_test.go):

<pre>
    go test -bench=. -benchtime=5s github.com/losalamos/goop
</pre>

On my computer, I get results like the following (reformatted for
clarity):

<pre>
    BenchmarkNativeFNV1             2000000000                 4.59 ns/op
    BenchmarkNativeFNV1Closure      2000000000                 4.59 ns/op
    BenchmarkGoopFNV1                100000000                70.5  ns/op
    BenchmarkMoreGoopFNV1             10000000               794    ns/op
    BenchmarkEvenMoreGoopFNV1          5000000              2517    ns/op
</pre>

See [`goop_test.go`](http://github.com/losalamos/goop/blob/master/goop_test.go) for the complete source code for those benchmarks.  Basically,

* `BenchmarkNativeFNV1` is native (i.e., non-Goop) Go code for computing a 64-bit [FNV-1 hash](http://isthe.com/chongo/tech/comp/fnv/) on a sequence of 0xFF bytes.  Each iteration ("op" in the performance results) comprises a nullary function call, a multiplication by a large prime number, and an exclusive or with an 0xFF byte.

* `BenchmarkNativeFNV1Closure` is the same but instead of calling an ordinary function each iteration, it invokes a locally defined closure.

* `BenchmarkGoopFNV1` defines a Goop object that contains a single data field (the current hash value) and no methods.  Each iteration performs one `Get` and one `Set` on the Goop object.

* `BenchmarkMoreGoopFNV1` replaces the hash function with an object method.  Hence, each iteration performs one `Get`, one `Set`, and one `Call` on the Goop object.

* `BenchmarkEvenMoreGoopFNV1` adds support for type-dependent dispatch to the hash-function method.  Although only one type signature is defined, Goop has to confirm at run time that the provided arguments do in fact match that signature.

Another way to interpret the data shown above is that, on my computer at least, a function closure is essentially free; `Get` and `Set` cost approximately 33&nbsp;ns apiece; a `Call` of a nullary function costs about 724&nbsp;ns; and type-dependent dispatch costs an additional 1723&nbsp;ns.

How does Goop compare to various scripting languages?  Not well, at least for `BenchmarkMoreGoopFNV1` and its equivalents in other languages.  The following table shows the cost in nanoseconds of an individual `BenchmarkMoreGoopFNV1` operation (a function call, a read of a data field, a 64-bit multiply, an 8-bit exclusive&nbsp;or, and a write to a data field):

<table style="border-collapse: collapse; margin-left: auto; margin-right: auto">
  <tr>
    <th style="text-align: left; border-top: solid medium; border-bottom: solid thin">Language</th>
    <th style="text-align: right; border-top: solid medium; border-bottom: solid thin">Run time (ns/op)</th>
  </tr>
  <tr>
    <td>[Incr Tcl] 8.6.0</td>
    <td style="text-align: right">24490</td>
  </tr>
  <tr style="background-color: yellow">
    <td>Go 1.3.1 + Goop</td>
    <td style="text-align: right">794</td>
  </tr>
  <tr>
    <td>JavaScript 1.7 (Rhino 1.7R4-2)</td>
    <td style="text-align: right">682</td>
  </tr>
  <tr>
    <td>Perl 5.18.2</td>
    <td style="text-align: right">622</td>
  </tr>
  <tr>
    <td>Python 2.7.6</td>
    <td style="text-align: right">604</td>
  </tr>
  <tr>
    <td>Python 3.4.0</td>
    <td style="text-align: right">567</td>
  </tr>
  <tr>
    <td>Ruby 1.9.3.4</td>
    <td style="text-align: right">513</td>
  </tr>
  <tr>
    <td style="border-bottom: solid medium">Python 2.7.3 + PyPy </td>
    <td style="border-bottom: solid medium; text-align: right">206</td>
  </tr>
</table>

In short, you'll want to do most of your coding in native Go and use Goop only when your application requires the extra flexibility that Goop provides.  Then, you should cache as many object members as possible in Go variables to reduce the number of `Get` and `Set` calls.

License
-------

Goop is provided under a BSD-ish license with a "modifications must be indicated" clause.  See <http://github.com/losalamos/goop/blob/master/LICENSE> for the full text.

Los Alamos National Security, LLC (LANS) owns the copyright to Goop, a component of the LANL Go Suite (identified internally as LA-CC-11-056).

Author
------

Scott Pakin, <pakin@lanl.gov>
