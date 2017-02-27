{
alert("javascript loaded.");

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
  console.log(folderId);
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
  console.log(clickedFolder);
  if (clickedFolder.expanded == true) {
    console.log(clickedDiv);
    document.getElementById('' + clickedDiv + '-content').innerHTML = "";
    clickedFolder.expanded = false;
  }
  else {
    for (let referenceId of clickedFolder.references) {
      reference = getFolderById(referenceId);
      $('#' + clickedDiv + '-content').append(
        '<div id="' + reference.id + '"> ' + reference.name +
          '<div id="' +reference.id + '-content"> </div> ' +
        '</div>');
      $('#' + clickedDiv + '-content').click(clickFolder);

    }
    clickedFolder.expanded = true;
  }
};

var folderButton1 = document.getElementById("1");
folderButton1.addEventListener('click', clickFolder);
}