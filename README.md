# xmonad-log

xmonad-log is a DBus monitoring solution that can easily be used to display
xmonad in a statusbar like [polybar](https://github.com/jaagr/polybar),
[lemonbar](https://github.com/LemonBoy/bar) and similar.

## Installation

xmonad-log is written in Go with one dependency:
[dbus](https://github.com/godbus/dbus). [Binary
packages](https://github.com/xintron/xmonad-log/releases) are available.

### Building from source

This package has been tested with Go 1.7 and above.

To build from source:
 1. Clone this repository into `$GOPATH/src/github.com/xintron/xmonad-log`.
 2. Build it within the directory with `go build`.

This should leave a `xmonad-log` binary in the directory. Move this to an
appropriate directory in your `$PATH`.

## Configure xmonad

To configure xmonad to send log events over DBus the haskell
[dbus](http://hackage.haskell.org/package/dbus) package is required. Once
installed the following can be added to your `.xmonad/xmonad.hs` configuration
to add DBus support.

```haskell
import XMonad
import XMonad.Hooks.DynamicLog

import qualified DBus as D
import qualified DBus.Client as D
import qualified Codec.Binary.UTF8.String as UTF8

main :: IO ()
main = do
    dbus <- D.connectSession
    -- Request access to the DBus name
    D.requestName dbus (D.busName_ "org.xmonad.Log")
        [D.nameAllowReplacement, D.nameReplaceExisting, D.nameDoNotQueue]

    xmonad $ def { logHook = dynamicLogWithPP (myLogHook dbus) }

-- Override the PP values as you would otherwise, adding colors etc depending
-- on  the statusbar used
myLogHook :: D.Client -> PP
myLogHook dbus = def { ppOutput = dbusOutput dbus }

-- Emit a DBus signal on log updates
dbusOutput :: D.Client -> String -> IO ()
dbusOutput dbus str = do
    let signal = (D.signal objectPath interfaceName memberName) {
            D.signalBody = [D.toVariant $ UTF8.decodeString str]
        }
    D.emit dbus signal
  where
    objectPath = D.objectPath_ "/org/xmonad/Log"
    interfaceName = D.interfaceName_ "org.xmonad.Log"
    memberName = D.memberName_ "Update"
```

View [this
xmonad-config](https://github.com/xintron/configs/blob/22a33b41587c180172392f80318883921c543053/.xmonad/lib/Config.hs#L199)
for a fully working polybar example using statusbar coloring.
