(describe str
          (== (str 0) "0")
          (== (str 1.4) "1.4")
          (== (str "1.0") "1.0")
          (== (str "hi") "hi")
          (== (str 'a) "a")
          (== (str '(1 2)) "(1 2)")
          (== (str '(1 . 2)) "(1 . 2)")
          (== (str "abc" 1 "-" 34.2 '(a b c)) "abc1-34.2(a b c)"))