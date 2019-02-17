name := "Monkey";
age := 1;
inspirations := ["Scheme", "Lisp", "JavaScript", "Clojure"];
book := {
  "title": "Writing A Compiler In Go",
  "author": "Thorsten Ball",
  "prequel": "Writing An Interpreter In Go"
};

printBookName := fn(book) {
    title := book["title"];
    author := book["author"];
    print(author + " - " + title);
};

printBookName(book);

fibonacci := fn(x) {
  if (x == 0) {
    return 0
  } else {
    if (x == 1) {
      return 1;
    } else {
      return fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

map := fn(arr, f) {
  iter := fn(arr, accumulated) {
    if (len(arr) == 0) {
      return accumulated
    } else {
      return iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };

  return iter(arr, []);
};

numbers := [1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2];
map(numbers, fibonacci);
