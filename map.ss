(define kcons
  (lambda (x y k)
    (k (cons x y))))

(define fmap
  (lambda (f l)
    (if (null? l)
        '()
        (cons (f (car l))
              (fmap f
                    (cdr l))))))

(define kmap
  (lambda (f l k)
    (if (null? l)
        (k '())
        (f (car l)
           (lambda (r1)
             (kmap f (cdr l)
                   (lambda (r2)
                     (k (cons r1 r2)))))))))

(define kinit
  (lambda (x) x))



'fmap

(fmap 'a '())
(fmap (lambda (x) (cons 'x x)) '(a b c))

'kmap

(kmap 'a '() kinit)
(kmap (lambda (x k) (k (cons 'x x))) '(a b c) kinit)

