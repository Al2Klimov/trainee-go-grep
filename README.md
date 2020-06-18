# trainee-grep

## About

**trainee-grep** is a from scratch developed program which outputs a pattern that is found in files you specify. There are
various parameters you can use to filter or individualize the output of your to "grepped" file.
Seen like this it's an imitation of the original grep command line command.

**Note:** You cannot use grep without at least the **-e** parameter which defines the pattern you search for.
## Parameters

| Parameter    | Usage | 
| ------------- |-------------|
| **-e**      | Define a pattern for matching.
| **-F**      | Use the pattern not as a regular expression but as a  fixed string. 
| **-w**  | Use pattern that only matches words.    
| **-x**     | Use pattern that only matches whole lines.
| **-i**     | Ignore case distinctions.
| **-v**   | Use pattern as non-matching lines. 
| **-q**      | Suppress all normal output.
| **-m**   | Stop printing after NUM selected lines.
| **-n**   | Print the line number.
| **-r**   | Search the pattern in a directory.

## Usage Examples

In the following you'll see a quick example of how to use the grep command. Here is what our text file looks like:

toGrep.txt:

                zebra
                zebra
                dog
                zebra
                cat
                
toGrep2.txt:

                dog
                cat
                cat
                mouse
                zebra

Now we can grep our pattern with `./grep.linux-amd64 -e 'zebra' < toGrep.txt`  

The output is:

                zebra
                zebra
                zebra
                
If you want to search for a pattern in more than one file use `./grep.linux-amd64 -e 'zebra' toGrep.txt toGrep2.txt `

The output is:

                toGrep.txt:zebra
                toGrep.txt:zebra
                toGrep.txt:zebra
                toGrep2.txt:zebra
