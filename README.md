# ahist

Simple history of search requests in **Acme**

## Indroduction

**ahist** tracks of search requests in any **Acme**'s window and keep them in a separate window.

## Using

Run **ahist** from **Acme**'s window and all requests made by `Look` command or `B3` mouse button will be tracked in a separate window with `+History` in a tag.
Tracking can be stopped by `-ahist` command from the window at any time.
**ahist** can be used with [atag](https://github.com/santucco/atag)

## Bugs

* **Acme** does not reflect immediately that a window is modified. So an appearance of `Put` command can be late a bit.
* **ahist** conflicts with other **Acme**'s applications that open `event` file, because **Acme** sends events only to a last connected client of `event` file.
