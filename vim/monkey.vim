" Vim Syntax File
" Language:     monkey
" Creator:      James Mills, prologic at shortcircuit dot net dot au
" Last Change:  31st January 2019

if version < 600
    syntax clear
elseif exists("b:current_syntax")
    finish
endif

syntax case match

syntax keyword xType true false null

syntax keyword xKeyword let fn if else return while

syntax keyword xFunction len input print first last rest push pop exit assert

syntax keyword xOperator == != < > !
syntax keyword xOperator + - * / =

syntax region xString start=/"/ skip=/\\./ end=/"/

syntax region xComment start='#' end='$' keepend
syntax region xComment start='//' end='$' keepend

highlight link xType Type
highlight link xKeyword Keyword
highlight link xFunction Function
highlight link xString String
 highlight link xComment Comment
highlight link xOperator Operator
highlight Operator ctermfg=5

let b:current_syntax = "monkey"
