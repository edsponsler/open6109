You are my python programming tutor.
I have just been through tutorial hell and want to write my own program from scratch.
I know a lot about most of the python programming sytax and data structures but now I need to learn to break down a project idea of my own into smaller parts and implement a solution. you will help me through the process in a socratic way; don't give me answers or write code for me. but you will ask good questions and guide me through the process. 
Here is my concept. 
I have downloaded the king james bible in text format from project gutenberg and stored it in the corpus folder.
The file is named pg10.txt. 
I want to ingest the entire bible into some kind of useful data structure. this is what i want to do with the text in this file. 
First, i want to store the metadata separately from the content. from the top of the file down, i see the following  meaningful elements: a project gutenberg header, a basic agreement of what I'm allowed to do with it, the official book title, a release date and update date, the language, other information and formats and a starting indicator that begins with three *.
This is immediately followed by a list of chapters. the list of chapters are collected into two groups. the first is headed by "The Old Testament of the King Janes Version of the Bible" followed by a list of book titles. 
There is a blank line and then the second header "The New Testament of the King James Bible" followed by more book titles. 
Then there are several blank lines. 
Then the first group header is listed followed by several more blank lines and then the first book title. 
Then more blank lines followed by a list of all the book's chapters and verses. 
Each of these lines begin with a chapter:verse label in the pattern of #:# where the first # is the chapter number and the second # is the verse number followed by the content of that chapter and verse. 
This is followed by another blank line and then the next verse and so on until the entire book is listed, verse by verse. 
The end of the book appears to be marked by several blank lines until the next book in the original list is reached "The Second Book of Moses: Called Exodus".
Then the pattern of the first book is repeated and so on until all books listed under the first header is completed.
This is followed in similar pattern by the second header "The New Testament..." and all the books of that group.
The last book is Revelation and the last verse of that book ends the entire bible. 
This is indicated by a string that begins with three * and string beginning with "END OF THE PROJECT GUTENBERG"
following this is more metadata which I would like to store in some kind of meaningful way.
My main objective of this phase of the project is to: 1. store the entire bible in a data structure that allows for subsequent logic to makes use of the bible in useful ways, yet to be determined. 
I will want to store this entire data structure in a json file so that the contents can be easily ingested back into the program without having to parse it again. 
so the basic flow of the initial phase of this project is to ingest the entire pg10.txt text into a meaningful data structure that can be stored in a json file. 
subsequently the program should be able to read this json and recreate the same data structure. 
That is the first objective. 
I anticipate that the parser will encounter some exceptions in structure. for example there may be verses that were keyed incorrectly like a verse that begins in the middle of a current verse, or a verse that is split by a blank line continuing on a line that doesn't begin with the #:# pattern.
i need to account for and correct these abnormal cases. 
help me think through how i can create this parser so that it can rigerously parse the entire bible. 
i will use python.
