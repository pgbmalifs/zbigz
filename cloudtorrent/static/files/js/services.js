/* globals app,window */

app.service('api', function($rootScope, $http, reqerr) {
  window.http = $http;
  var request = function(action, data) {
    var url = "/api/"+action;
    $rootScope.apiing = true;
    return $http.post(url, data).error(reqerr).finally(function() {
      $rootScope.apiing = false;
    });
  };
  var api = {};
  var actions = ["configure","magnet","url","torrent","file"];
  actions.forEach(function(action) {
    api[action] = request.bind(null, action);
  });
  return api;
});

app.service('search', function($rootScope, $http, reqerr) {
  return {
    all: function(provider, query, page) {
      var params = {query:query};
      if(page !== undefined) params.page = page;
      $rootScope.searching = true;
      var req = $http.get("/search/"+provider, { params: params });
      req.error(reqerr);
      req.finally(function() {
        $rootScope.searching = false;
      });
      return req;
    },
    one: function(provider, path) {
      var opts = { params: { path:path } };
      $rootScope.searching = true;
      var req = $http.get("/search/"+provider+"-item", opts);
      req.error(reqerr);
      req.finally(function() {
        $rootScope.searching = false;
      });
      return req;
    }
  };
});

app.service('storage', function() {
  return window.localStorage || {};
});

app.service('hash', function() {
  return function(obj) {
    return JSON.stringify(obj, function(k,v) {
      return k.charAt(0) === "$" ? undefined : v;
    });
  };
});

app.service('reqerr', function() {
  return function(err, status) {
    alert(err.error || err);
    console.error("request error '%s' (%s)", err, status);
  };
});

app.service('hexid', function() {
  var N = 12;
  return function() {
    return Array(N+1).join((Math.random().toString(16)+'00000000000000000').slice(2, 18)).slice(0, N);
  };
});

app.service('bytes', function() {
  var scale = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  return function(n) {
    var i = 0;
    var s = scale[i];
    if (typeof n !== 'number') {
      return "-";
    }
    while (n > 1000) {
      s = scale[++i] || 'x10^' + (i * 3);
      n = Math.round(n / 100) / 10;
    }
    return "" + n + " " + s;
  };
});
