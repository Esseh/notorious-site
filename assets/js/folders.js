{

var parentFolder1 = {
    name: "Parent Folder 1",
    id: 1,
    isFolder: true,
    references: [3, 4, 7],
    expanded: false
};

var parentFolder2 = {
    name: "Parent Folder 2",
    id: 2,
    isFolder: true,
    references: [testNote],
    expanded: false
};

var childFolder1 = {
    name: "Child Folder 1",
    id: 3,
    isFolder: true,
    references: [5, 6, 7],
    expanded: false
};

var childFolder2 = {
    name: "Child Folder 2",
    id: 4,
    isFolder: true,
    references: [7],
    expanded: false
};

var childOfChildFolder1 = {
    name: "Child-Child Folder 1",
    id: 5,
    isFolder: true,
    references: [7],
    expanded: false
};

var childOfChildFolder2 = {
    name: "Child-Child Folder 2",
    id: 6,
    isFolder: true,
    references: [7],
    expanded: false
};

var testNote = {
    name: "Test Note",
    id: 7,
    isFolder: false,
    references: [],
    expanded: false
};

var userFolders = [parentFolder1,
                   parentFolder2,
                   childFolder1,
                   childFolder2,
                   childOfChildFolder1,
                   childOfChildFolder2,
                   testNote];

var getFolderById = function(folderId) {
  for (let folder of userFolders) {
    if (folderId == folder.id) {
      return folder;
    }
  }
  alert("folderID not found");
};

var clickFolder = function (event) {
  console.log(event.target.id);
  var clickedDiv = event.target.id;
  var clickedFolder = getFolderById(clickedDiv);
  if (clickedFolder.expanded == true) {
    document.getElementById('' + clickedDiv + '-content').innerHTML = "";
    $(document.getElementById('' + clickedDiv + '-content')).unbind();
    $('#' + clickedDiv + '-content').removeClass("open-content");
    $('#' + clickedDiv + '-content').addClass("content");
    $('#' + clickedDiv).removeClass("open-folder");
    $('#' + clickedDiv).addClass("folder");
    clickedFolder.expanded = false;
  }
  else {
    for (let referenceId of clickedFolder.references) {
      reference = getFolderById(referenceId);
      $('#' + clickedDiv + '-content').append(
        '<div id="' + reference.id + '" class="folder"> ' + reference.name + '</div> </div>' +
          '<div id="' +reference.id + '-content" class="content"> </div> ');
      $('#' + reference.id).unbind();
      $('#' + reference.id).click(clickFolder);

    }
    $('#' + clickedDiv + '-content').removeClass("content");
    $('#' + clickedDiv + '-content').addClass("open-content");
    $('#' + clickedDiv).removeClass("folder");
    $('#' + clickedDiv).addClass("open-folder");
    console.log(clickedDiv);
    clickedFolder.expanded = true;
  }
};

var addFolder = function (parentString, nameString) {
  $.post('/folder/api/newfolder', { parent: parentString, name: nameString }, function (data) {
    console.log(data);
  });
};

var removeFolder = function (parentString) {
  $.post('/folder/api/deletefolder', { parent: parentString }, function (data) {
    console.log(data);
  });
};

var initializeRoot = function () {
  $.post('/folder/api/initializeroot', {}, function (data) {
    console.log(data);
  });
}

document.getElementById("test-button").onclick = function () {
  initializeRoot();
};

$(document.getElementById("1")).unbind();
$(document.getElementById("1")).click(clickFolder);
}