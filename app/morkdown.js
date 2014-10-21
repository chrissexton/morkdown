"use strict";

angular.module('MorkdownApp', [
	'angular-markdown-editable',
	'rt.debounce',
	'cfp.hotkeys'
])
.service('morkdownBackend', function($http) {
	return {
		get: function(path) {
			return $http.get('/files/'+path);
		},
		save: function(path, text) {
			return $http.put('/files/'+path, text);
		},
		delete: function(path) {
			return $http.delete('/files/'+path);
		},
		list: function(path) {
			return $http.post('/files/'+path, {});
		}
	};
})
.controller('TextEditorController', function($scope, debounce, hotkeys, morkdownBackend) {
	var newDoc = "# New Document\nThis is a new document for you!"

	$scope.list = function(path) {
		return morkdownBackend.list(path);
	}

	$scope.new = function() {
		$scope.text = newDoc;
		$scope.file = "";
		$scope.modified = false;
	};

	$scope.open = function(force) {
		if ($scope.modified && !force) {
			return;
		}
		morkdownBackend.get($scope.file)
			.then(function(res) {
				$scope.text = res.data;
				$scope.modified = false;
			})
			.catch(function(err) {
				if (err.status === 404)
					$scope.text = newDoc;
				else
					$scope.text = err.data;
				$scope.modified = false;
			});
	};

	$scope.onFileChange = debounce(1000, $scope.open);


	$scope.save = function() {
		morkdownBackend.save($scope.file, $scope.text)
			.then(function(data) {
				$scope.modified = false;
			});
	};

	hotkeys.add({
		combo: "alt+o",
		description: "Load file",
		allowIn: ['INPUT', 'SELECT', 'TEXTAREA'],
		callback: $scope.open.bind($scope, true)
	});

	hotkeys.add({
		combo: "alt+n",
		description: "New file",
		allowIn: ['INPUT', 'SELECT', 'TEXTAREA'],
		callback: $scope.new.bind($scope)
	});

	hotkeys.add({
		combo: "alt+s",
		description: "Save file",
		allowIn: ['INPUT', 'SELECT', 'TEXTAREA'],
		callback: $scope.save.bind($scope)
	});

	$scope.new();

});
