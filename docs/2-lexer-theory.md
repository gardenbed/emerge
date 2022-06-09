# Theory of Lexers

## Regular Expressions

Regular expressions are the standard notation for specifying lexeme patterns.
Each regular expression `r` denotes a language `L(r)`,
which is also defined recursively from the languages denoted by `r` subexpressions.
A language that can be defined by a regular expression is called a **regular set**.

## Finite Automata

Finite automata are recognizers. They simply say *yes* or *no* about each possible input string.
Finite automata come in two flavors:

  - **Nondeterministic Finite Automata (NFA)**
  - **Deterministic Finite Automata (DFA)**

We can represent an NFA or DFA by a *transition graph*, where the nodes are states and the labeled edges represent the transition function.
There is an edge labeled a from state `s` to state `t` if and only if `t` is one of the next states for state `s` and input `a`.

### Nondeterministic Finite Automata

A nondeterministic finite automaton (NFA) consists of:

  1. A finite set of states `S`.
  2. A set of input symbols `Σ`, the *input alphabet*. We assume that the empty string `ε` is never a member of `Σ`.
  3. A *transition function* that gives, for each state, and for each symbol in `Σ ∪ {ε}` a set of next states.
  4. A state <code>s<sub>0</sub> ∈ S</code> that is distinguished as the *start state* or *initial state*.
  5. A set of states `F ⊆ S` that is distinguished as the *accepting states* or *final states*.

Note that:

  - The same symbol can label edges from one state to several difierent states.
  - An edge may be labeled by the empty string `ε` instead of, or in addition to, symbols from the input alphabet.

We can also represent an NFA by a *transition table*, whose rows correspond to states, and whose columns correspond to the input symbols and `ε`.
The entry for a given state and input is the value of the transition function applied to those arguments.

### Deterministic Finite Automata

A deterministic finite automaton (DFA) is a special case of an NFA where:

  - There are no moves on input `ε`.
  - For each state `s` and input symbol `a`, there is exactly one edge out of `s` labeled `a`.

Every NFA (and every regular expression) can be converted to a DFA accepting the same language.

## Core Algorithms

### Simulating a DFA

```
# INPUT:  A DFA D with start state s0, accepting states F, and transition function move.
#         An input string x terminated by an end-of-file character EOF.
# OUTPUT: Answer "yes" if D accepts x; "no" otherwise.

s = s0
c = nextChar()

while (c != EOF) {
  s = move(s, c)
  c = nextChar()
}

if (s is in F) {
  return true
}
return false
```

### The Subset Construction

```
# INPUT:  An NFA N.
# OUTPUT: A DFA D accepting the same language as N.
# METHOD: The algorithm constructs a transition table Dtran for D.
#         Each state of D is a set of NFA states.
#         Dtran is constructed such so D will simulate, in parallel, all possible moves N can make on a given input string.

initially, ε-closure({s0}) is the only state in Dstates, and it is unmarked
while (there is an unmarked state T in Dstates) {
  mark T;
  for (each input symbol a) {
    U = ε-closure(move(T, a))
    if (U is not in Dstates) {
      add U as an unmarked state to Dstates
    }
    Dtran[T, a] = U
  }
}

# Computing ε-closure(T)
push all states of T onto stack
initialize ε-closure(T) to T
while (stack is not empty) {
  pop t, the top element, off stack
  for (each state u with an edge from t to u labeled ε) {
    if (u is not in ε-closure(T)) {
      add u to ε-closure(T)
      push u onto stack
    }
  }
}
```

### Simulating an NFA

```
# INPUT:  An NFA N with start state s0, accepting states F, and transition function move.
#         An input string x terminated by an end-of-file character EOF.
# OUTPUT: Answer "yes" if N accepts x; "no" otherwise.
# METHOD: The algorithm keeps a set of current states S,
#         those that are reached from s0 following a path labeled by the inputs read so far.
#         If c is the next input character, read by the function nextChar(),
#         then we first compute move(S, c) and then close that set using ε-closure().

S = ε-closure(s0)
c = nextChar()

while (c != EOF) {
  S = ε-closure(move(S, c))
  c = nextChar()
}

if (S ∩ F != ∅) {
  return true
}
return false
```

To implement the above algorithm in an efficient way, we need the following data structures:

  - **Two stacks**, each of which holds a set of NFA states.
    - One of these stacks, `oldStates`, holds the current set of states.
    - The second, `newStates`, holds the next set of states.
  - A boolean array `alreadyOn`, indexed by the NFA states, to indicate which states are in `newStates`.
    While the array and stack hold the same information,
    it is much faster to interrogate `alreadyOn[s]` than to search for state `s` on the stack newStates.
  - A two-dimensional array `move[s, a]` holding the transition table of the NFA.
    The entries in this table, which are sets of states, are represented by linked lists.

```
# Adding a new state s, which is known not to be on newStates
func addState(s) {
  push s onto newStates
  alreadyOn[s] = TRUE
  for (t on move[s, ε]) {
    if (!alreadyOn[t]) {
      addState(t)
    }
  }
}

# Implementation of ε-closure(move(S, c))
for (s on oldStates) {
  for (t on move[s, c]) {
    if (!alreadyOn[t]) {
      addState(t)
    }
  }
  pop s from oldStates
}

for (s on newStates) {
  pop s from newStates
  push s onto oldStates
  alreadyOn[s] = FALSE
}
```

The *running time* of the above algorithm is `O(k(n + m))`.

### Converting Regular Expression to NFA

The **McNaughton-Yamada-Thompson** algorithm is used for converting a regular expression to an NFA.
The algorithm is **syntax-directed**, in the sense that it works recursively up the parse tree for the regular expression.

```
# INPUT:     A regular expression r over alphabet Σ.
# OUTPUT:    An NFA N accepting L(r).
# METHOD:    Begin by parsing r into its constituent subexpressions.
#            The rules for constructing an NFA consist of basis rules for handling subexpressions with no operators,
#            and inductive rules for constructing larger NFA's from the NFA's for the immediate subexpressions of a given expression.
#
# INDUCTION: Suppose N(s) and N(t) are NFA's for regular expressions s and t, respectively.
#
# r = s|t:   i and f are new states, the start and accepting states of N(r), respectively.
#            There are ε-transitions from i to the start states of N(s) and N(t),
#            and each of their accepting states have ε-transitions to the accepting state f.
#            The accepting states of N(s) and N(t) are not accepting in N(r) anymore.
#            N(r) accepts L(s) ∩ L(t), which is the same as L(r).
#
# r = st:    The start state of N(s) becomes the start state of N(r),
#            and the accepting state of N(t) is the only accepting state of N(r).
#            The accepting state of N(s) and the start state of N(t) are merged into a single state,
#            with all the transitions in or out of either state.
#            N(r) accepts exactly L(s)L(t), and is a correct NFA for r = st.
#
# r = s:     i and f are new states, the start state and lone accepting state of N(r).
#            To get from i to f, we can either follow the introduced path labeled,
#            which takes care of the one string in L(s)⁰, or we can go to the start state of N(s),
#            through that NFA, then from its accepting state back to its start state zero or more times.
#            These options allow N(r) to accept all the strings in L(s)¹, L(s)², and so on,
#            so the entire set of strings accepted by N(r) is L(s*).
#
# r = (s):   L(r) = L(s), and we can use the NFA N(s) as N(r).
#
```

## Implementation

One question that may arise is whether or not we should construct and simulate a DFA or an NFA.
Apparently, it is faster to have a DFA to simulate than an NFA.
In principle, the number of DFA states does not influence the running time of the DFA simulation algorithm.
We may favor NFA over DFA because of the fact that the *subset construction* can exponentiate the number of states in the worst case.

If the string-processor is going to be used many times (i.e., lexer) then any cost of converting to a DFA is worthwhile.
However, in other string-processing applications (i.e., grep), where the user specifies a regular expression each time,
it may be more eficient to skip the step of constructing a DFA and simulate the NFA directly.

| **Automaton** | **Construction** | **Simulation** |
|----|----|----|
| NFA | <code>O(\|r\|)</code> | <code>O(\|r\| × \|x\|)</code> |
| DFA (average case) | <code>O(\|r\|<sup>3</sup>)</code> |  <code>O(\|x\|)</code> |
| DFA (worst case) | <code>O(\|r\|<sup>2</sup>2<sup>\|r\|</sup>)</code> | <code>O(\|x\|)</code> |

### Lexer Architecture

```
                    ┌────────────┐       ┌──────────────────┐
    EBNF            │ Lexer      │       │ Transition Table │
 Description ──────►│ Generator  ├──────►├──────────────────┤
                    │ (Compiler) │       │ Actions          │
                    └────────────┘       └─────────▲────────┘
                                                   │
                                                   │
                                         ┌─────────▼────────┐
                                         │    Automaton     │
                                         │    Simulator     │
                                         └────┬───────┬─────┘
                                              │       │
                                  lexemeBegin │       │ forward
  Input Buffer                                │       │
 ┌────────────────────────────────────────────▼───────▼─────┐
 │                                             lexeme       │
 └──────────────────────────────────────────────────────────┘
```

### Input Buffering

Buffering techniques should be used to reduce the amount of overhead required to process a single input character.

A lexical analyzer may need to read ahead some characters before it can determine the next token.
A general approach to reading ahead on the input involves
maintaining an input buffer from which the lexer can read and push back characters.
A pointer keeps track of the portion of the input that has been analyzed;
pushing back a character is implemented by moving back the pointer.

The most common scheme involves two buffers that are alternately reloaded.
Each buffer is of the same size `N` (`N` is usually the size of a disk block,`4096` bytes).
If fewer than `N` characters remain in the input, then a special character (*EOF*) marks the end of the input.

Two pointers to the input are maintained:
  1. Pointer `lexemeBegin`, marks the beginning of the current lexeme, whose extent we are attempting to determine.
  2. Pointer `forward` scans ahead until a pattern match is found.

Once the next lexeme is determined, `forward` is set to the character at its right end.
Then, after the lexeme is recorded, `lexemeBegin` is set to the character immediately after the lexeme just found.

We can combine the buffer-end test with the test for the current character
if we extend each buffer to hold a sentinel character at the end.
The sentinel is a special character that cannot be part of the source program, and a natural choice is *EOF*.

Note that as long as we never need to look so far ahead of the actual lexeme
that the sum of the lexeme's length plus the distance we look ahead is greater than `N`,
we shall never overwrite the lexeme in the buffer before determining it.
**If character strings can be very long, extending over many lines,
then we could face the possibility that a lexeme is longer than `N`.**

### Reserved Words vs. Identifiers

There are two ways to recognize reserved words:

  1. Initialize the reserved words in the symbol table.
     A field of the symbol table entry indicates that these strings are never ordinary identifiers,
     and tells which token they represent.
  2. Create separate finite automata for each keyword.
     If we adopt this approach, then we must prioritize the tokens
     so that the reserved-word tokens are recognized in preference to id, when the lexeme matches both patterns.

### Lexical Analysis using DFAs

To construct an automaton for a lexical analyzer,

  1. Convert each regular expression pattern in the input description
     to an NFA using [this](#converting-regular-expression-to-nfa) algorithm.
  2. Combine all the NFA's into one by introducing a new start state with ε-transitions
     to each of the start states of the NFA's N<sub>i</sub> for pattern p<sub>i</sub>.
     (at this point, we can simulate the NFA directly or proceed to the next step).
  3. Convert the the NFA into an equivalent DFA using the [subset construction](#the-subset-construction) algorithm.
  4. Simulate the DFA until at some point there is no next state
     (or strictly speaking, the next state is `∅`, the dead state corresponding to the empty set of NFA states).
  5. Within each DFA state, if there are one or more accepting NFA states, determine the first pattern
     whose accepting state is represented, and make that pattern the output of the DFA state.

#### The Lookahead Operator

When converting the lookahead operator in a pattern like `r1/r2` to an NFA,
we use the `ε` for `/`, and we will not look for a `/` on the input.
If the NFA recognizes a prefix `xy` of the input buffer as matching this regular expression,
the end of the lexeme is not where the NFA entered its accepting state.
Rather, the end occurs when the NFA enters a state `s` such that

  1. `s` has an ε-transition on the `/`.
  2. There is a path from the start state of the NFA to state `s` that spells out `x`.
  3. There is a path from state `s` to the accepting state that spells out `y`.
  4. `x` is as long as possible for any `xy` satisfying conditions 1-3.

If there is only one ε-transition state on the imaginary `/` in the NFA,
then the end of the lexeme occurs when this state is entered for the last time.

#### Dead States in DFAs

We must know when there is no longer any possibility of recognizing a longer lexeme.
Thus, we always omit transitions to the *dead state* and eliminate the *dead state* itself.

If we construct a DFA from a regular expression using the
[McNaughton-Yamada-Thompson](#converting-regular-expression-to-nfa) and [subset construction](#the-subset-construction) algorithms,
then the DFA will not have any states besides `∅` that cannot lead to an accepting state.

### Optimizing DFA-Based Lexers

  1. The first algorithm constructs a DFA directly from a regular expression, without constructing an intermediate NFA.
     The resulting DFA also may have fewer states than the DFA constructed via an NFA.
  2. The second algorithm minimizes the number of states of any DFA, by combining states that have the same future behavior.
     The algorithm itself is quite efficient, running in time `O(nlogn)`, where `n` is the number of states of the DFA.
  3. The third algorithm produces a more compact representation of transition tables than the standard, two-dimensional table.

#### Converting Regular Expression Directly to DFA

##### Important States

A state of an NFA is **important** if it has a non-ε out-transition.
The *subset construction* algorithm uses only the important states in a set `T`
when it computes `ε-closure(move(T, a))`, the set of states reachable from `T` on input `a`.
The set of states `move(s, a)` is non-empty only if state `s` is important.

During the *subset construction*, two sets of NFA states can be treated as the same set if they:

  1. Have the same important states, and
  2. Either both have accepting states or neither does.

When an NFA is constructed from the [McNaughton-Yamada-Thompson](#converting-regular-expression-to-nfa) algorithm,
each important state corresponds to a particular operand in the regular expression.

By concatenating a unique right end-marker `~` to a regular expression `r`,
we give the accepting state for `r` a transition on `~`,
making it an important state of the NFA for `(r)~`.
By using the augmented regular expression `(r)~`,
we can forget about accepting states as the subset construction proceeds.
When the construction is complete, any state with a transition on `~` must be an accepting state.

The important states of the NFA correspond directly to
the positions in the regular expression that hold symbols of the alphabet.

##### Functions Computed from Syntax Tree

To construct a DFA directly from a regular expression, we construct its syntax tree
and then compute four functions: **nullable**, **firstPos**, **lastPos**, and **followPos**,

  - `nullable(n)` is `true` for a syntax-tree node `n` if and only if the subexpression represented by `n` has `ε` in its language
    (the empty string `ε` is in the language of the subexpression rooted at `n`).
  - `firstPos(n)` is the set of positions in the subtree rooted at `n` that
     correspond to the first symbol of at least one string in the language of the subexpression rooted at `n`.
  - `lastPos(n)` is the set of positions in the subtree rooted at `n` that
    correspond to the last symbol of at least one string in the language of the subexpression rooted at `n`.
  - `followPos(p)`, for a position `p`, is the set of positions `q` in the entire syntax tree
    such that there is some string <code>x = a<sub>1</sub>a<sub>2</sub>...a<sub>n</sub></code> in `L((r)#)`
    such that for some `i`, there is a way to explain the membership of `x` in `L((r)#)`
    by matching <code>a<sub>i</sub></code> to position `p` of the syntax tree and <code>a<sub>i</sub>+1</code> to position `q`.

| **Node `n`** | **`nullable(n)`** | **`firstPos(n)`** | **`lastPos(n)`** |
|----------|----------|----------|----------|
| Leaf labeled `ε` | `true` | `∅` | `∅` |
| Leaf with position `i` | `false` | `{i}` | `{i}` |
| Star <code>n*</code> | `true` | <code>firstPos(n)</code> | <code>lastPos(n)</code> |
| Alt <code>n<sub>1</sub>\|n<sub>2</sub></code> | <code>nullable(n<sub>1</sub>) \|\| nullable(n<sub>2</sub>)</code> | <code>firstPos(n<sub>1</sub>) ∪ firstPos(n<sub>2</sub>)</code> | <code>lastPos(n<sub>1</sub>) ∪ lastPos(n<sub>2</sub>)</code> |
| Concat <code>n<sub>1</sub>n<sub>2</sub></code> | <code>nullable(n<sub>1</sub>) && nullable(n<sub>2</sub>)</code> | <code>if nullable(n<sub>1</sub>)</br>&nbsp;&nbsp;firstPos(n<sub>1</sub>) ∪ firstPos(n<sub>2</sub>)</br>else</br>&nbsp;&nbsp;firstPos(n<sub>1</sub>)</code> | <code>if nullable(n<sub>2</sub>)</br>&nbsp;&nbsp;lastPos(n<sub>1</sub>) ∪ lastPos(n<sub>2</sub>)</br>else</br>&nbsp;&nbsp;lastPos(n<sub>2</sub>)</code> |

For computing the `followPos(p)` function,
there are only two ways that a position of a regular expression can be made to follow another.

  1. If `n` is a *concat*-node with left child <code>n<sub>1</sub></code> and right child <code>n<sub>2</sub></code>,
     then for every position `i` in <code>lastPos(n<sub>1</sub>)</code>,
     all positions in <code>firstPos(n<sub>2</sub>)</code> are in `followPos(i)`.
  2. If `n` is a *star*-node, and `i` is a position in <code>lastPos(n<sub>1</sub>)</code>,
     then all positions in <code>firstPos(n<sub>1</sub>)</code> are in `followPos(i)`.

We can represent the function `followPos` by creating a directed graph
with a node for each position and an edge from position `i` to position `j` if and only if `j` is in `followPos(i)`.
The graph for `followPos` would become an NFA without ε-transitions for the underlying regular expression if we:

  1. Make all positions in `firstPos` of the root be initial states.
  2. Label each edge from `i` to `j` by the symbol at position `i`.
  3. Make the position associated with end-marker `~` be the only accepting state.

##### Algorithm

```
# INPUT:  A regular expression r.
# OUTPUT: A DFA D that recognizes L(r).
# METHOD:
#   1. Construct a syntax tree T from the augmented regular expression (r)~.
#   2. Compute nullable, firstPos, lastPos, and followPos for T.
#   3. Construct Dstates, the set of states of DFA D, and Dtran, the transition function for D.
#      The states of D are sets of positions in T.
#      Initially, each state is "unmarked," and a state becomes "marked" just before we consider its out-transitions.
#      The start state of D is firstPos(n0), where node n0 is the root of T.
#      The accepting states are those containing the position for the end-marker symbol ~.

initialize Dstates to contain only the unmarked state firstPos(n0), where n0 is the root of syntax tree T for (r)~
while (there is an unmarked state S in Dstates) {
  mark S
  for (each input symbol a) {
    let U be the union of followPos(p) for all p in S that correspond to a
    if (U is not in Dstates) {
      add U as an unmarked state to Dstates
    }
    Dtran[S, a] = U
  }
}
```

#### Minimizing The Number of DFA States

If we implement a lexer as a DFA, we would generally prefer a DFA with as few states as possible,
since each state requires entries in the transition table that describes the lexer.

There is always a unique minimum state DFA for any regular language.
This minimum-state DFA can be constructed from any DFA for the same language by grouping sets of equivalent states.

String `x` distinguishes state `s` from state `t`
if exactly one of the states reached from `s` and `t` by following the path with label `x` is an accepting state.
State `s` is distinguishable from state `t` if there is some string that distinguishes them.

##### Algorithm

```
# INPUT:  A DFA D with set of states S, input alphabet Σ, start state s0, and set of accepting states F.
# OUTPUT: A DFA D' accepting the same language as D and having as few states as possible.
# METHOD:
#   1. Start with an initial partition Π with two groups, F and S - F, the accepting and non-accepting states of D.
#   2. Apply the following procedure to construct a new partition Πnew.

initially, let Πnew = Π
for (each group G of Π) {
  partition G into subgroups such that two states s and t are in the same subgroup if and only if
  for all input symbols a, states s and t have transitions on a to states in the same group of Π
  // at worst, a state will be in a subgroup by itself
  replace G in Πnew by the set of all subgroups formed
}

#   3. If Πnew = Π, let Πfinal = Π and continue with step 4.
#      Otherwise, repeat step 2 with Πnew in place of Π.
#   4. Choose one state in each group of Πfinal as the representative for that group.
#      The representatives will be the states of the minimum-state DFA D'.
#      The other components of D' are constructed as follows.
#
#        (a) The start state of D' is the representative of the group containing the start state of D.
#        (b) The accepting states of D' are the representatives of those groups that contain an accepting state of D.
#            Each group contains either only accepting states, or only non-accepting states.
#        (c) Let s be the representative of some group G of Πfinal, and let the transition of D from s on input a be to state t.
#            Let r be the representative of t's group H. Then in D', there is a transition from s to r on input a.
```

##### Eliminating the Dead State

The minimization algorithm sometimes produces a DFA with one dead state.
One that is not accepting and transfers to itself on each input symbol.
This state is technically needed, because a DFA must have a transition from every state on every symbol.
We often want to know when there is no longer any possibility of acceptance, so we can establish that the proper lexeme has already been seen.
We may wish to eliminate the dead state and use an automaton that is missing some transitions.
This automaton has one fewer state than the minimum-state DFA,
but is strictly speaking not a DFA, because of the missing transitions to the dead state.

#### Trading Time for Space in DFA Simulation

The simplest and fastest way to represent the transition function of a DFA is a 2D table indexed by states and characters.
In situations where memory resource is limited (embedded devices),
there are many methods that can be used to compact the transition table.

There is a data structure that allows us to combine the speed of array access with the compression of lists with defaults.

```
    default   base               next    check
   ┌────────┬────────┐        ┌────────┬────────┐
   │        │        │        │        │        │
   │        │        │        │        │        │
   │        │        │        │        │        │
   │        │        │        │        │        │
   │        │        │        │        │        │
   ├────────┼────────┤        │        │        │
 s │   q    │   ─────┼───────►│        │        │
   ├────────┼────────┤      ▲ │        │        │
   │        │        │      │ │        │        │
   │        │        │      │ │        │        │
   │        │        │      a │        │        │
   │        │        │      │ │        │        │
   │        │        │      │ │        │        │
   │        │        │      ▼ ├────────┼────────┤
   │        │        │        │   r    │   t    │
   │        │        │        ├────────┼────────┤
   │        │        │        │        │        │
   │        │        │        │        │        │
   │        │        │        │        │        │
   └────────┴────────┘        └────────┴────────┘
```

We may think of this structure as four arrays.
The `base` array is used to determine the base location of the entries for state `s`,
which are located in the `next` and `check` arrays.
The `default` array is used to determine an alternative `base` location
if the `check` array tells us the one given by `base[s]` is invalid.
The `nextState` function is defined as follows:

```
int nextState(s, a) {
  if (check[base[s] + a] == s) {
    return next[base[s] + a]
  }
  return nextState(default[s], a)
}
```

The intended use of the above data structure is to
make the `next-check` arrays short by taking advantage of the similarities among states.

While we may not be able to choose base values so that no `next-check` entries remain unused,
experience has shown that the simple strategy of assigning base values to states in turn, and
assigning each `base[s]` value the lowest integer so that the special entries for state `s` are not previously occupied
utilizes little more space than the minimum possible.
