# Vim/Vi Keyboard Shortcut Cheat Sheet

> [!NOTE]
> This cheat sheet is designed for experienced users returning to Vim. It focuses on efficiency, navigation, search-and-replace, splits, and other intermediate/advanced shortcuts that are easy to forget.

## 1. Advanced Navigation

| Shortcut | Description |
| :--- | :--- |
| `H` / `M` / `L` | Move cursor to **H**igh, **M**iddle, or **L**ow of screen |
| `zt` / `zz` / `zb` | Redraw screen with current line at the **t**op, **z**enter (center), or **b**ottom |
| `Ctrl + d` / `Ctrl + u` | Scroll **d**own or **u**p half-page |
| `Ctrl + f` / `Ctrl + b` | Scroll **f**orward or **b**ackward full page |
| `Ctrl + e` / `Ctrl + y` | Scroll screen down or up by 1 line (keeps cursor position if possible) |
| `w` / `b` | Move forward / backward by word (alphanumeric/underscore) |
| `W` / `B` | Move forward / backward by WORD (separated by whitespace only) |
| `e` / `ge` | Move to end of word / end of previous word |
| `E` / `gE` | Move to end of WORD / end of previous WORD |
| `0` / `^` / `$` | Move to start of line / first non-blank character / end of line |
| `g_` | Move to last non-blank character of the line |
| `}` / `{` | Move forward / backward by paragraph (blank lines) |
| `)` / `(` | Move forward / backward by sentence |
| `%` | Jump to matching parenthesis, bracket, or brace |
| `gd` | Go to local definition of symbol under cursor |
| `gg` / `G` | Jump to first line / last line of file |
| `:{num}` or `{num}G` | Jump to line number `{num}` |
| `Ctrl + o` / `Ctrl + i` | Jump to **o**lder / **i**nner (newer) cursor position in jump list |

## 2. Marks and Jumps

Marks allow you to save cursor positions and jump back to them.

| Shortcut | Description |
| :--- | :--- |
| `m{a-z}` | Set mark `{a-z}` at current cursor position (file-local) |
| `m{A-Z}` | Set global mark `{A-Z}` (works across files) |
| `'{a-z}` | Jump to the start of the line containing mark `{a-z}` |
| ``{a-z}` | Jump to the exact character position of mark `{a-z}` |
| `''` / ` `` ` | Jump back to the position before the last jump / line before the last jump |
| `:marks` | List all current marks |

## 3. Editing, Replacing & Substitution

### Text Objects (The "Verb + Modifier + Object" pattern)
*Verbs: `d` (delete), `c` (change), `y` (yank/copy), `v` (select)*  
*Modifiers: `i` (inner/inside), `a` (around/including)*

| Shortcut | Action | Description |
| :--- | :--- | :--- |
| `ciw` / `caw` | Change Word | Change **i**nner word / **a**round word (including trailing space) |
| `ci"` / `ca"` | Change in Quotes | Change inside double quotes / around double quotes (includes quotes) |
| `ci(` or `cib` | Change in Parens | Change inside parentheses / around parentheses |
| `ci{` or `ciB` | Change in Braces | Change inside curly braces / around curly braces |
| `cit` | Change in Tag | Change inside XML/HTML tags |
| `diw` / `yiw` | Delete / Yank | Delete / copy inner word |

### Basic Replacements
| Shortcut | Description |
| :--- | :--- |
| `r{char}` | Replace single character under cursor with `{char}` |
| `R` | Enter Replace (overwrite) mode |
| `s` | Delete character under cursor and enter Insert mode |
| `S` or `cc` | Delete line and enter Insert mode |
| `~` | Toggle case of character under cursor |
| `g~{motion}` | Toggle case of text covered by `{motion}` (e.g., `g~w` toggles word case) |
| `gu{motion}` | Lowercase text covered by `{motion}` |
| `gU{motion}` | Uppercase text covered by `{motion}` |

### Search and Replace (Substitution command `:s`)
| Command | Description |
| :--- | :--- |
| `:s/old/new` | Replace first occurrence of `old` with `new` in current line |
| `:s/old/new/g` | Replace all occurrences of `old` with `new` in current line |
| `:%s/old/new/g` | Replace all occurrences of `old` with `new` in the entire file |
| `:%s/old/new/gc` | Replace all with confirmation prompt (`y`/`n`/`a`/`q`) |
| `:10,20s/old/new/g` | Replace all occurrences of `old` with `new` between lines 10 and 20 |
| `:'<,'>s/old/new/g` | Replace all occurrences in the current Visual selection |

## 4. Copy, Cut, Paste & Registers

Vim has multiple clipboards called **registers**.

| Shortcut | Description |
| :--- | :--- |
| `y` / `yy` / `Y` | Yank (copy) selection / yank line / yank to end of line |
| `d` / `dd` / `D` | Delete (cut) selection / delete line / delete to end of line |
| `p` / `P` | Paste after cursor / paste before cursor |
| `"{reg}y` | Yank into register `{reg}` (e.g., `"ayw` copies word to register `a`) |
| `"{reg}p` | Paste from register `{reg}` (e.g., `"ap` pastes contents of register `a`) |
| `"+y` / `"+p` | Yank to / paste from system clipboard (requires `+clipboard` compilation) |
| `"0p` | Paste from the *last yank* register (avoids pasting deleted text) |
| `:reg` / `:registers` | View contents of all registers |

## 5. Visual Modes (Selection)

| Shortcut | Description |
| :--- | :--- |
| `v` | Start visual mode (character-based selection) |
| `V` | Start visual line mode (line-based selection) |
| `Ctrl + v` | Start visual block mode (column/matrix selection) |
| `o` | Move cursor to other end of visual selection |
| `I` / `A` | Insert/Append text in Visual Block mode (type text, then press `Esc` to apply to all lines) |
| `>` / `<` | Indent / outdent selected lines |

## 6. Multi-Window and Tab Management

### Split Windows
| Command / Shortcut | Description |
| :--- | :--- |
| `:sp {file}` / `Ctrl + w s` | Split window horizontally |
| `:vsp {file}` / `Ctrl + w v` | Split window vertically |
| `Ctrl + w w` | Cycle cursor through open windows |
| `Ctrl + w h/j/k/l` | Move cursor left, down, up, or right to adjacent window |
| `Ctrl + w H/J/K/L` | Move current window to far left, bottom, top, or right |
| `Ctrl + w =` | Resize all split windows to equal dimensions |
| `Ctrl + w _` / `Ctrl + w \|` | Maximize window height / maximize window width |
| `Ctrl + w c` or `:q` | Close current window |
| `Ctrl + w o` or `:only` | Close all windows except the current one |

### Tabs
| Command / Shortcut | Description |
| :--- | :--- |
| `:tabnew {file}` | Open `{file}` in a new tab |
| `:tabn` / `gt` | Go to next tab |
| `:tabp` / `gT` | Go to previous tab |
| `:tabfirst` / `:tablast` | Go to first tab / last tab |
| `:tabonly` | Close all other tabs |

## 7. Macros and Automation

Macros allow you to record and replay sequences of keystrokes.

| Shortcut | Description |
| :--- | :--- |
| `q{a-z}` | Start recording keystrokes into register `{a-z}` |
| `q` | Stop recording macro |
| `@{a-z}` | Execute macro recorded in register `{a-z}` |
| `@@` | Repeat the last executed macro |
| `5@{a-z}` | Execute the macro 5 times |
| `.` | Repeat last change/edit command (incredibly powerful with search/motions) |

## 8. Miscellaneous Helpers

| Shortcut / Command | Description |
| :--- | :--- |
| `J` | Join current line with the line below (adds a space) |
| `Ctrl + a` / `Ctrl + x` | Increment / decrement number under the cursor |
| `Ctrl + p` / `Ctrl + n` | Autocomplete word (backward / forward search) in Insert mode |
| `:%normal {commands}` | Run normal-mode `{commands}` on every line of the file (e.g. `:%normal i// ` comments out a file) |
| `:r {file}` | Read content of `{file}` and insert it below cursor |
| `:r !{cmd}` | Run shell `{cmd}` and insert its standard output below cursor |
