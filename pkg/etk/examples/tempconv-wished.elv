use etk
use str

var temp-conv-data = [
  &celsius=[
    &other=fahrenheit
    &conveter={|c| + 32 (* 9/5 $c)}
  ]
  &fahrenheit=[
    &other=celsius
    &converter={|f| * 5/9 (- $f 32)}
  ]
]

# Need decorator here
fn temp-conv {|this:| [etk:comp]
  this:subcomp celsius = (etk:textbox [&prompt=(styled 'Celsius: ')])
  this:subcomp fahrenheit = (etk:textbox [&prompt=(styled 'Fahrenheit: ')])
  this:state focus = celsius
  pprint (dissoc $this:state state-dump) |
    str:replace "\t" "" (slurp) |
    put (styled "\n State dump: \n" inverse)(one) |
    this:subcomp state-dump = (etk:textbox [&prompt=(one)])

  put [
    # Need tag here
    &layout=[%etk:vbox celsius fahrenheit state-dump]
    &focus=$this:state[focus]
    &handler={|e| [returnable]
      val other = $temp-conv-data[$focus][other]
      if (==s $e Tab) {
        set this:state[focus] = $other
        return
      }

      if (not (this:propagate $focus $e)) {
        etk:not-handled
        return
      }
    }
  ]
}

edit:push-addon (etk:adapt-to-widget $temp-conv~)
