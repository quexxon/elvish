use etk
use str

var temp-conv = (etk:comp {|this:|
  this:state focus = (num 0)

  etk:vbox [
    (this:subcomp celsius $etk:textbox [&prompt=(styled 'Celsius: ')])
    (this:subcomp fahrenheit $etk:textbox [&prompt=(styled 'Fahrenheit: ')])

    (this:subcomp state-dump $etk:textbox [&prompt=
      (pprint (dissoc $this:state state-dump) |
       str:replace "\t" "" (slurp) |
       put (styled "\n State dump: \n" inverse)(one))])
  ] [
    &focus=$this:state[focus]
    &handler={|e es|
      if (==s $es Tab) {
        set this:state[focus] = (- 1 $this:state[focus])
      } else {
        if (== $this:state[focus] 0) {
          if (etk:handle $this:subcomp[celsius] $e) {
            try {
              var f = (+ 32 (* 9/5 $this:state[celsius][buffer][content]) | printf '%.2f' (one))
              set this:state[fahrenheit][buffer] = (etk:text-buffer $f (count $f))
            } catch {
            }
          } else {
            # Can we do better than this? 🤔
            put $false
          }
        } else {
          if (etk:handle $this:subcomp[fahrenheit] $e) {
            try {
              var c = (* 5/9 (- $this:state[fahrenheit][buffer][content] 32) | printf '%.2f' (one))
              set this:state[celsius][buffer] = (etk:text-buffer $c (count $c))
            } catch {
            }
          } else {
            put $false
          }
        }
      }
    }
  ]
})

edit:push-addon (etk:adapt-to-widget $temp-conv)
