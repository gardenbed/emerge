# Definitions

In this document, we define two languages using formal mathematical notations.
The first language is the [Regular Expression](https://en.wikipedia.org/wiki/Regular_expression),
and the second language is the [Extended Backus-Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form).

To distinguish between *terminal* and *non-terminal* symbols in this document,
we use lowercase letters or quotation marks for terminal symbols and uppercase letters for non-terminal symbols.

We use regular expressions for defining patterns for *tokens* (*terminal symbols*)
in the EBNF language itself as well as any input languages desribed in EBNF.
We need a regular expression parser for building a lexer for any EBNF parser.

## Regular Expression

Regular expression, as the name suggests, is a regular language or a Type-3 language.
Since Type-2 languages are a superset of Type-3 languages,
we can define the regular expression language using the EBNF notation as well
(we will define the EBNF notation formally in the next section).

### Grammar

{% raw %}
```
regex            = [ "^" ] expr
expr             = subexpr [ "|" expr ]
subexpr          = {{ subexpr_item }}
subexpr_item     = anchor | group | match
anchor           = "$"
group            = "(" expr ")" [ quantifier ]
match            = match_item [ quantifier ]
match_item       = any_char | single_char | char_class | ascii_char_class | char_group
char_group       = "[" [ "^" ] {{ char_group_item }} "]"
char_group_item  = ascii_char_class | char_class | char_range | single_char
char_range       = char_in_range "-" char_in_range
char_in_range    = unicode_char | ascii_char | char
quantifier       = repetition [ "?" ]
repetition       = "?" | "*" | "+" | range
range            = "{" num [ upper_bound ] "}"
upper_bound      = "," [ num ]
any_char         = "."
single_char      = unicode_char | ascii_char | escaped_char | unescaped_char
char_class       = "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
ascii_char_class = "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"

ascii_char       = "\x" hex_digit{2}
unicode_char     = "\x" hex_digit{4,8}
escaped_char     = "\" ( "\" | "|" | "." | "?" | "*" | "+" | "(" | ")" | "[" | "]" | "{" | "}" | "$" )
unescaped_char   = # all characters excluding the escaped ones
char             = # all characters

num              = {{ digit }}
letters          = {{ letter }}
digit            = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
hex_digit        = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" | "A" | "B" | "C" | "D" | "E" | "F"
letter           = "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z"
                 | "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z"
```
{% endraw %}

## Extended Backus-Naur Form

The formal language we use for describing [context-free grammars](https://en.wikipedia.org/wiki/Context-free_grammar)
is a subset of [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) notation.

### Alphabet

```
"  /  =  |  (  )  [  ]  {  }
0  1  2  3  4  5  6  7  8  9
A  B  C  D  E  F  G  H  I  J  K  L  M  N  O  P  Q  R  S  T  U  V  W  X  Y  Z
a  b  c  d  e  f  g  h  i  j  k  l  m  n  o  p  q  r  s  t  u  v  w  x  y  z
!  #  $  %  &  '  *  +  ,  -  .  :  ;  <  >  ?  @  \  ^  _  `  ~
```

### Tokens

{% raw %}
| **Token** | **Lexeme**  | **Description**|
|-----------|-------------|----------------|
| `DEF`     | `"="`  | Symbol for *rule* definition |
| `SEMI`    | `";"`  | Symbol for end of *rule* definition |
| `ALT`     | `"\|"` | Symbol for *alternation* |
| `LPAREN`  | `"("`  | Start symbol for *grouping* |
| `RPAREN`  | `")"`  | End symbol for *grouping* |
| `LBRACK`  | `"["`  | Start symbol for *optional* |
| `RBRACK`  | `"]"`  | End symbol for *optional* |
| `LBRACE`  | `"{"`  | Start symbol for *repetition (Kleene Star)* |
| `RBRACE`  | `"}"`  | End symbol for *repetition (Kleene Star)* |
| `LLBRACE` | `"{{"` | Start symbol for *repetition (Kleene Plus)* |
| `RRBRACE` | `"}}"` | End symbol for *repetition (Kleene Plus)* |
| `LANGLE`  | `"<"`  | Star symbol for *rule* reference |
| `RANGLE`  | `">"`  | End symbol for *rule* reference |
| `PREDEF`  | `/\$[A-Z][0-9A-Z_]*/` | predefined *token* definitions |
| `LASSOC`  | `"@left"`  | keyword for specifying left-associative terminals |
| `RASSOC`  | `"@right"` | keyword for specifying right-associative terminals |
| `NOASSOC` | `"@none""` | keyword for specifying non-associative terminals |
| `GRAMMER` | `"grammar"` | Keyword for declaring grammar name |
| `IDENT`   | `/[a-z][0-9a-z_]*/` | Regex for grammar name and *non-terminal* symbols |
| `TOKEN`   | `/[A-Z][0-9A-Z_]*/` | Regex for declaring and referencing *terminal* symbols |
| `STRING`  | `/"([\x21\x23-\x5B\x5D-\x7E]\|\\[\x21-\x7E])+"/` | Regex for defining *string* patterns |
| `REGEX`   | `/\/([\x20-\x2E\x30-\x5B\x5D-\x7E]\|\\[\x20-\x7E])*\//` | Regex for defining *regular expression* patterns |
{% endraw %}

### Grammar

{% raw %}
```
grammar   = name {decl};
name      = "grammar" IDENT [";"];
decl      = token [";"] | directive [";"] | rule ";";
token     = TOKEN "=" (STRING | REGEX | PREDEF);
directive = ("@left" | "@right" | "@none") {{term | "<" rule ">"}};
rule      = lhs "=" [rhs];
lhs       = nonterm;
rhs       = rhs rhs | "(" rhs ")" | "[" rhs "]" | "{" rhs "}" | "{{" rhs "}}" | rhs "|" rhs | rhs "|" | nonterm | term;
nonterm   = IDENT;
term      = TOKEN | STRING;
```
{% endraw %}

  - Tokens can be defined **explicitly** using *token declarations*
    or **implicitly** by using *string literals* in rule definitions.
  - A semicolon at the end of *grammar name*, *token definitions*, or *directives* is optional.
  - The rule `rule = lhs "=" [rhs]` allows the definition of *empty productions*, such as `A → ε`.
  - The rule `rhs = rhs "|"` permits trailing empty alternatives in production rules,
    enabling definitions like `A → B | C | ε`.
  - The addition of the `";"` token helps distinguish between successive rule definitions,
    eliminating ambiguity in the grammar. Other ways to achieve this include:
      - Using newlines to separate rules and indentation to break a rule definition into multiple lines.
      - Introducing a special keyword at the beginning of each rule to explicitly mark its start.
      - Enclosing each rule with special delimiter tokens to clearly indicate its boundaries.
      - Modifying the rule `rhs → rhs rhs` to `rhs → rhs "," rhs`
        to prevent unintended chaining of the next rule's head into the current rule’s body.
  - The `"<"` and "`>`" tokens allow multiple rules to be referenced in the same line when
    specifying *associativity* and *precedence*, ensuring they are not mistakenly mixed with actual rule definitions.

### Precedence and Associativity

The grammar above, as it stands, is ambiguous.
Below are the associativity and precedence rules for an LR parser to resolve the conflicts.

{% raw %}
```
@left  <rhs = rhs rhs>
@left  "(" "[" "{" "{{" IDENT TOKEN STRING
@right "|"
@none  "="
@none  "@left" "@right" "@none"
```

  1. The production rule `rhs = rhs rhs` (concatenating two `rhs`)
     has the highest precedence and is *left-associative*.
  2. The next highest precedence is assigned to the terminals
     `"("`, `"["`, `"{"`, `"{{"`, `IDENT`, `TOKEN`, and `STRING`, all of which are *left-associative*.
  3. The `"|"` terminal is assigned the next level of precedence and is *right-associative*.
  4. Finally, the `"="` terminal has the lowest precedence and is *non-associative*.
{% endraw %}
