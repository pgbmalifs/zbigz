/* globals app,window */

app.controller("ConfigController", function($scope, $rootScope, storage, api, hash, hexid) {
  $rootScope.config = $scope;
  //inputs is a copy of configurations
  $scope.inputs = {};
  $scope.saved = true;
  $scope.$watch("inputs", function() {
    $scope.saved = $scope.cfgsHash === hash($scope.inputs);
  }, true);
  //on module changes, extract configurations
  $rootScope.$watch("state.Modules", function(modules) {
    var cfgs = {};
    Object.values(modules || {}).forEach(function(m) {
      cfgs[m.TypeID] = m.Config;
    });
    var cfgsHash = hash(cfgs);
    if($scope.cfgsHash === cfgsHash) {
      return;
    }
    $scope.cfgsHash = cfgsHash;
    $scope.cfgs = cfgs;
    $scope.inputs = angular.copy(cfgs);
    $scope.saved = true;
  }, true);

  $scope.submitConfig = function() {
    api.configure($scope.inputs);
  };
  $scope.add = function(type) {
    var typeid = type + ":" + hexid();
    $scope.inputs[typeid] = {};
  };
  $scope.list = function(type) {
    var ms = [];
    for(var typeid in $scope.inputs) {
      if(typeid.indexOf(type) === 0) {
        ms.push($scope.inputs[typeid]);
      }
    }
    return ms;
  };
});
