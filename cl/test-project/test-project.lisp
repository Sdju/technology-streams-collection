;;;; test-project.lisp

;;; (in-package #:test-project)

;;; (css-lite:css (("body") (:width "70%")))

;;; (css-lite:css (("body") (:width "70%")))
;;; (ql:quickload "com.inuoe.jzon")
;;; (ql:quickload 'css-lite)

(defparameter *design-system* "")

(defparameter *design-system-file* (open "C:\\pr\\my\\streams\\cl\\test-project\\assets\\design-system.json" :if-does-not-exist nil))

(when *design-system-file*
            (loop :for line = (read-line *design-system-file* nil)
                :while line :do (setf *design-system* (concatenate 'string *design-system* line)))
            (close *design-system-file*))

(defmacro dohash ((table &key key val) &body body) 
    (cond
     ((and (null key) (null val)) (error "WrongSyntax"))
     ((null key) `(loop :for ,val :being :the :hash-values :in ,table :do ,@body))
     ((null val) `(loop :for ,key :being :the :hash-keys :in ,table :do ,@body))
     (t `(loop :for ,key :being :the :hash-keys :in ,table
           :using (hash-value ,val)
           :do ,@body))))

(defmethod print-json-object ((a float) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "~a~%" a))

(defmethod print-json-object ((a number) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "~d~%" a))

(defmethod print-json-object ((a (eql t)) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "TRUE~%"))

(defmethod print-json-object ((a (eql nil)) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "FALSE~%"))

(defmethod print-json-object ((a symbol) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "~s~%" a))

(defmethod print-json-object ((a string) &optional (padding "") (direction t)) 
    (declare (ignore padding))
    (format direction "~s~%" a))

(defun make-ostring ()
    (make-array '(0) :element-type 'base-char :fill-pointer 0 :adjustable t))

(defmethod print-json-object ((a vector) &optional (padding "") (direction t)) 
    (if (eql direction t)
        (progn
          (format t "[~%")
          (let ((newPadding (concatenate 'string padding "  ")))
            (loop :for i :across a do 
                    (format t newPadding)
                    (print-json-object i newPadding)))
          (format t "~a]~%" padding))
        (let (
              (result (format nil "[~%"))
              (newPadding (concatenate 'string padding "  ")))
          (loop :for i :across a do 
                    (setf result (concatenate 'string result 
                                       (format nil newPadding)
                                       (print-json-object i newPadding nil))))
                    
          (setf result (concatenate 'string result (format nil "~a]~%" padding)))
          result)))

(defmethod print-json-object ((a hash-table) &optional (padding "") (direction t)) 
    (if (eql direction t)
        (progn
          (format t "{~%")
          (let ((newPadding (concatenate 'string padding "  ")))
            (dohash (a :key key :val value)
                    (format t "~a[~s] => " newPadding key)
                    (print-json-object value newPadding)))
          (format t "~a}~%" padding))
        (let (
              (result (make-ostring))
              (newPadding (concatenate 'string padding "  ")))
          (with-output-to-string (s result)
            (format s "{~%")
            (dohash (a :key key :val value)
                (format s "~a[~s] => ~a" newPadding key (print-json-object value newPadding nil)))
            (format s "~a]~%" padding))
          result)))

(defun make-class (selector properties)
    (let ((result (make-ostring)))
        (with-output-to-string (s result)
            (format s "~a {~%" selector)
            (dohash (properties :key key :val value)
                    (format s "  ~a: ~a;~%" key value))
            (format s "}~%"))
        result))

(let* ((result (make-ostring))
       (design-system (com.inuoe.jzon:parse *design-system*))
       (design-system-colors (gethash "colors" design-system))
       (design-system-spaces (gethash "spaces" design-system))
       (all-properties (make-hash-table :test 'equal))
       (result-file (open "C:\\pr\\my\\streams\\cl\\test-project\\assets\\design-system.css" :direction :output :if-exists :supersede)))
    (with-output-to-string (s result)
        (dohash (design-system-colors :key key :val value)    
            (let* ((property-name (format nil "--color-~a" key))
                   (property-name-usage (format nil "var(~a, ~a)" property-name value))
                   (c-properties (make-hash-table :test 'equal))
                   (bg-properties (make-hash-table :test 'equal)))
                   
                (setf (gethash property-name all-properties) value)
                (setf (gethash "color" c-properties) property-name-usage)
                (setf (gethash "background-color" bg-properties) property-name-usage)
                
                (format s "~a~%" (make-class
                                  (format nil ".c-~a" key)
                                  c-properties))
                (format s "~a~%" (make-class
                                  (format nil ".bc-~a" key)
                                  bg-properties))))
                
        (dohash (design-system-spaces :key key :val value)
            (let* ((property-name (format nil "--spaces-~a" key))
                   (property-name-usage (format nil "var(~a, ~a)" property-name value))
                   (p-properties (make-hash-table :test 'equal))
                   (m-properties (make-hash-table :test 'equal)))
                   
                (setf (gethash property-name all-properties) value)
                (setf (gethash "padding" p-properties) property-name-usage)
                (setf (gethash "margin" m-properties) property-name-usage)
                
                (format s "~a~%" (make-class
                                  (format nil ".p-~a" key)
                                  p-properties))
                (format s "~a~%" (make-class
                                  (format nil ".m-~a" key)
                                  m-properties))))
        (format s "~a~%" (make-class ":root" all-properties)))
            
    (format result-file "~a" result)
    ;;; (format result-file "~a" (css-lite:css (("body") (:width "70%"))))
    (close result-file))
