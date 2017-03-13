{

var clickFolder = function (event) {
  var clickedDiv = event.target.id;
  openFolder(clickedDiv);
  $(document.getElementById('' + clickedDiv)).unbind();
  $(document.getElementById('' + clickedDiv)).click(clickOpenFolder);
};

var clickOpenFolder = function (event) {
  var clickedDiv = event.target.id;
    document.getElementById('' + clickedDiv + '-content').innerHTML = "";
    $(document.getElementById('' + clickedDiv + '-content')).unbind();
    $(document.getElementById('' + clickedDiv + '-content')).removeClass("open-content");
    $(document.getElementById('' + clickedDiv + '-content')).addClass("content");
    $(document.getElementById('' + clickedDiv)).removeClass("open-folder");
    $(document.getElementById('' + clickedDiv)).addClass("folder");
    $(document.getElementById('' + clickedDiv)).unbind();
    $(document.getElementById('' + clickedDiv)).click(clickFolder);
};

var clickAddFolder = function (event) {
  let baseFolder = event.target.value;
  console.log(baseFolder);
  let input = document.getElementById("" + baseFolder + "-user-input");
  console.log(input.value);
}

var clickRemoveFolder = function (event) {
  let baseFolder = event.target.value;
  console.log(baseFolder);
}

var openFolder = function (folderID) {
  let idString = '' + folderID;
  $.post('/folder/api/openfolder', { FolderID: idString }, function (data) {
    console.log(data);

    // var dataObj = $.parseJSON(data);
    // for (let referenceName of dataObj.folders) {
    // for (let noteName of dataObj.notes) {
    let mockFolders = ["test1", "test2"];
    let mockNotes = ["This is a good test note.", "This is an awesome test note."];

    $(document.getElementById("" + folderID + '-content')).append(
        '<button id="' + folderID + '-remove-folder" class="remove-folder" value="' + folderID + '"> Delete Folder </button>' +
        '<input id="' + folderID + '-user-input" class="user-input" type=text value="" placeholder=""/>' +
        '<button id="' + folderID + '-add-folder" class="add-folder" value="' + folderID + '"> New Folder </button>' +
        '<button id="' + folderID + '-add-note" class="add-note" value="' + folderID + '"> Add Note </button>');
    $(document.getElementById("" + folderID + "-add-folder")).unbind();
    $(document.getElementById("" + folderID + "-add-folder")).click(clickAddFolder);
    $(document.getElementById("" + folderID + "-remove-folder")).unbind()
    $(document.getElementById("" + folderID + "-remove-folder")).click(clickRemoveFolder);
    for (let referenceName of mockFolders) {
      let referenceId = "" + folderID + "/" + referenceName;
      $(document.getElementById("" + folderID + '-content')).append(
          '<div id="' + referenceId + '" class="folder"> ' + referenceName + '</div>' +
          // this ended with double </div>, assure no errors now that second div has been removed.
          '<div id="' +referenceId + '-content" class="content"> </div> ');
      $(document.getElementById(referenceId)).unbind();
      $(document.getElementById(referenceId)).click(clickFolder);
    }
    for (let noteName of mockNotes) {
      let noteId = "" + folderID + "/" + noteName;
      $(document.getElementById("" + folderID + '-content')).append(
          '<div class="note-container">' +
            '<div id="' + noteId + '-remove-note" class="remove-note" value="' + folderID + '"> X </div>' +
            '<div id="' + noteId + '" class="note"> ' + noteName + '</div>' +
          '</div>');
      $(document.getElementById(noteId + "-remove-note")).unbind();
      $(document.getElementById(noteId + "-remove-note")).click(clickRemoveNote);
      $(document.getElementById(noteId)).unbind();
      $(document.getElementById(noteId)).click(openNote);
    }
    $(document.getElementById('' + folderID + '-content')).removeClass("content");
    $(document.getElementById('' + folderID + '-content')).addClass("open-content");
    $(document.getElementById('' + folderID)).removeClass("folder");
    $(document.getElementById('' + folderID)).addClass("open-folder");
  });
};

var addFolder = function (parentId, folderName) {
  let parentIdString = "" + parentId;
  $.post('/folder/api/newfolder', { ParentID: parentIdString, FolderName: folderName }, function (data) {
    console.log(data);
  });
};

var removeFolder = function (parentId, folderName) {
  let parentIdString = "" + parentId;
  $.post('/folder/api/deletefolder', { ParentID: parentIdString, FolderName: folderName }, function (data) {
    console.log(data);
  });
};

var addNote = function (parentId, noteId) {
  let parentIdString = "" + parentId;
  let noteIdInt = parseInt(noteID, 10);
  $.post('/folder/api/addnote', { ParentID: parentIdString, NoteID: noteIdInt }, function (data) {
    console.log(data);
  });
}

var removeNote = function (parentId, noteId) {
  let parentIdString = "" + parentId;
  let noteIdInt = parseInt(noteID, 10);
  $.post('/folder/api/removenote', { ParentID: parentIdString, NoteID: noteIdInt }, function (data) {
    console.log(data);
  });
}

var initializeRoot = function () {
  $.post('/folder/api/initializeroot', { RootID: "5629499534213120" }, function (data) {
  });
  $(document.getElementsByClassName('root')).unbind();
  $(document.getElementsByClassName('root')).click(clickFolder);
};

var clickRemoveNote = function () {
  console.log('remove note called');
}

var openNote = function () {
  console.log('Open Note called');
  // will eventually navigate to note.
}

// ID of folders will me parentID+/+folderName+/+folderName+/+

// folders on returned data has no quotes around it so Json.parsing doesn't work.

$(document.getElementsByClassName('root')).unbind();
$(document.getElementsByClassName('root')).click(initializeRoot);
$(document.getElementsByClassName('root')).click(clickFolder);

}