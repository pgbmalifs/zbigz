

app.filter('keys', function() {
  return Object.keys;
});

app.filter('addspaces', function() {
  return function(s) {
    if(typeof s !== "string")
      return s;
    return s.replace(/([A-Z]+[a-z]*)/g, function(_, word) {
      return " " + word;
    }).replace(/^\ /, "");
  };
});

app.filter('filename', function() {
  return function(path) {
    return (/\/([^\/]+)$/).test(path) ? RegExp.$1 : path;
  };
});

app.filter('bytes', function(bytes) {
  return bytes;
});
