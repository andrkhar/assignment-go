set shell=/bin/bash

set encoding=utf-8
set fileencoding=utf-8

" Keeps 1000 items in the history.
set history=200

" Keep buffers in memmory
set hidden

" Interface
set ruler
set showcmd
set wildmenu
set number
"colorscheme gruvbox
"set background=dark

" e<CMD> lines before cursor
set scrolloff=5

" Search
set incsearch

" Prevents breaking words for line wrap
set lbr

" Keys mapping
map <leader>w :w!<CR>
map <leader>h :set hlsearch!<CR>

" Identation
set tabstop=4
set softtabstop=4
set shiftwidth=4
set expandtab
set autoindent
set smartindent

filetype plugin indent on
syntax on

" vim-go
let g:go_highlight_structs = 1 
let g:go_highlight_methods = 1
let g:go_highlight_functions = 1
let g:go_highlight_operators = 1
let g:go_highlight_build_constraints = 1

" vim: set ft=vim :
