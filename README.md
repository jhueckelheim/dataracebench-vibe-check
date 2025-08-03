# Vibe Check, based on DataRaceBench 1.4.0

This fork of [DataRaceBench](https://github.com/LLNL/dataracebench)
contains translations of the C and Fortran test cases to the Go
programming language. The translations were done by Claude 4 Sonnet
in July and August 2025. The translations are not perfect.
In particular, they

 - change the semantics of some test cases even in the absence of races;
 - introduce new races in some test cases; and
 - remove the races from some test cases that are expected to have races.

The goal of this translation was twofold:

 - Introduce and evaluate what we call a *translate-then-check* approach,
   where code from a source language A is translated to a target language
   B for the sole purpose of data race detection in the A code. We
   postulate that this can be helpful when data race detection tools are
   more mature for language B, or when language B is better suited for race
   detection.
 - Test whether the current AI agents can faithfully translate parallel
   programs. DataRaceBench is a useful evaluation tool for this, as it
   contains a number of small programs that exercise a large variety of
   parallel features in OpenMP, have been extensively studied,
   and at least about half of them are known to be free of data races.
