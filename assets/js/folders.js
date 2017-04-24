{

// NOTE: A folder's "value" holds is it's parent folder's ID. Value is used to store data for function calls on divs.

// NOTE: The root folder's id will be a number value associated with the user's profile.
// The ids of folders inside of the root follow a standard format.
// For example, say the rootID is "12345", a folder named "testFolder" located inside the root folder will have an id "12345/testfolder".
// Here is an example of a folder a few layers down, "12345/testFolder/childOfTestFolder/grandchildTestFolder".

//////////////////// GENERIC FUNCTIONS ////////////////////

// This function is used to refresh a folder's content to reflect changes made on the back end.
var refreshContent = function (folderId) {
  let divId = folderId;

  // Close folder content and set it to empty.
  document.getElementById('' + divId + '-content').innerHTML = "";
  $(document.getElementById('' + divId + '-content')).unbind();
  $(document.getElementById('' + divId + '-content')).removeClass("open-content");
  $(document.getElementById('' + divId + '-content')).addClass("content");
  $(document.getElementById('' + divId)).removeClass("open-folder");
  $(document.getElementById('' + divId)).addClass("folder");

  // Reopen the folder.
  openFolder(divId);
  $(document.getElementById('' + divId)).unbind();
  $(document.getElementById('' + divId)).click(clickOpenFolder);
}

//////////////////// ONCLICK FUNCTIONS ////////////////////

// This is the function called when a note is clicked. It will open the note in a new window.
var openNote = function (event) {
  window.open('/view/' + event.target.getAttribute("value"), "_blank");
};

// This is the function called when a closed folder is clicked.
// This will open the folder, close any open menus, open this folder's menu, and close all open folders that are not part of its parent chain.
var clickFolder = function (event) {

  // Open the folder.
  let clickedDiv = event.target.id;
  openFolder(clickedDiv);
  $(document.getElementById('' + clickedDiv)).unbind();
  $(document.getElementById('' + clickedDiv)).click(clickOpenFolder);

  // Close and empty all open folders that aren't in the parent chain.
  let parentId = document.getElementById(clickedDiv).getAttribute('value');
  let folderDivs = (document.getElementsByClassName("open-folder"));
    for (let folderDiv of folderDivs) {
      let folderID = folderDiv.getAttribute('id');
      if(!(('' + clickedDiv).indexOf('' + folderID + '/') != -1)) {
        if ((document.getElementById('' + folderID + '-content').innerHTML) != "") {
          document.getElementById('' + folderID + '-content').innerHTML = "";
          $(document.getElementById('' + folderID + '-content')).removeClass("open-content");
          $(document.getElementById('' + folderID + '-content')).addClass("content");
          $(document.getElementById('' + folderID)).removeClass("open-folder");
          $(document.getElementById('' + folderID)).addClass("folder");
          $(folderDiv).unbind();
          $(folderDiv).click(clickFolder);
        }
      }
    }
};

// This is the function called when an open folder is clicked.
// This will close the folder, all menus, and any folders inside of this one, the open the menu of it's parent folder.
var clickOpenFolder = function (event) {

  // Close and empty folder content.
  let clickedDiv = event.target.id;
  document.getElementById('' + clickedDiv + '-content').innerHTML = "";
  $(document.getElementById('' + clickedDiv + '-content')).unbind();
  $(document.getElementById('' + clickedDiv + '-content')).removeClass("open-content");
  $(document.getElementById('' + clickedDiv + '-content')).addClass("content");

  // Close the folder and change the on click function.
  $(document.getElementById('' + clickedDiv)).removeClass("open-folder");
  $(document.getElementById('' + clickedDiv)).addClass("folder");
  $(document.getElementById('' + clickedDiv)).unbind();
  $(document.getElementById('' + clickedDiv)).click(clickFolder);

  // Empty and close all menus.
  let menuDivs = (document.getElementsByClassName("open-menu"));
  for (let menuDiv of menuDivs) {
    menuDiv.innerHTML = "";
    $(menuDiv).removeClass('open-menu');
    $(menuDiv).addClass('menu');
  }

  // Close and empty all open folders that are not part of the parent chain.
  let parentId = document.getElementById(clickedDiv).getAttribute('value');
  let folderDivs = (document.getElementsByClassName("open-folder"));
    for (let folderDiv of folderDivs) {
      let folderID = folderDiv.getAttribute('id');
      if(!(('' + clickedDiv).indexOf('' + folderID + '/') != -1)) {
        if ((document.getElementById('' + folderID + '-content').innerHTML) != "") {
          document.getElementById('' + folderID + '-content').innerHTML = "";
          $(document.getElementById('' + folderID + '-content')).removeClass("open-content");
          $(document.getElementById('' + folderID + '-content')).addClass("content");
          $(document.getElementById('' + folderID)).removeClass("open-folder");
          $(document.getElementById('' + folderID)).addClass("folder");
          $(folderDiv).unbind();
          $(folderDiv).click(clickFolder);
        }
      }
    }

  //If the folder is not the root folder, then open the parent folder's menu and add content.
  if((document.getElementById('' + clickedDiv)).getAttribute('value') != "root") {
    $(document.getElementById("" + parentId + '-menu')).append(
        '<button id="' + parentId + '-remove-folder" class="remove-folder" value="' + parentId + '"> Delete Folder </button>' +
        '<button id="' + parentId + '-rename-folder" class="rename-folder" value="' + parentId + '"> Rename Folder </button>' +
        '<button id="' + parentId + '-add-folder" class="add-folder" value="' + parentId + '"> New Folder </button>' +
        '<button id="' + parentId + '-add-note" class="add-note" value="' + parentId + '"> Add Note </button>');
    $(document.getElementById('' + parentId + '-menu')).removeClass('menu');
    $(document.getElementById('' + parentId + '-menu')).addClass("open-menu");
    $(document.getElementById("" + parentId + "-add-folder")).unbind();
    $(document.getElementById("" + parentId + "-add-folder")).click(clickAddFolder);
    $(document.getElementById("" + parentId + "-add-note")).unbind();
    $(document.getElementById("" + parentId + "-add-note")).click(clickAddNote);
    $(document.getElementById("" + parentId + "-remove-folder")).unbind()
    $(document.getElementById("" + parentId + "-remove-folder")).click(clickRemoveFolder);
    $(document.getElementById("" + parentId + "-rename-folder")).unbind()
    $(document.getElementById("" + parentId + "-rename-folder")).click(clickRenameFolder);
  }

  //If the folder is the root folder, reopen the folder.
  else {
    openFolder(clickedDiv);
    $(document.getElementById('' + clickedDiv)).unbind();
    $(document.getElementById('' + clickedDiv)).click(clickOpenFolder);
  }
};

// This is the function called when Add Folder button is clicked. T
// This will create a prompt for an input value, and feed the input value and event data into the addFolder function.
var clickAddFolder = function (event) {
  let confirmFolderName = function() {
    let inputValue = document.getElementById('prompt-input').value;
    let baseFolder = event.target.value;
    addFolder(baseFolder, inputValue);
  }
  createPrompt('Enter the name of the folder you would like to create.', confirmFolderName);
}

var clickRenameFolder = function (event) {
  let confirmFolderName = function() {
    let inputValue = document.getElementById('prompt-input').value;
    let inputValue2 = document.getElementById('prompt-input2').value;
    let baseFolder = event.target.value;
    renameFolder(baseFolder, inputValue, inputValue2);
  }
  createDoublePrompt('Enter the name of the folder you would like to rename and the new name you would like to assign it.', confirmFolderName);
}

// This is the function called when Remove Folder button is clicked.
// This will take the click event info and feed it into the removeFolder function.
var clickRemoveFolder = function (event) {
  let confirmFolderName = function() {
    let inputValue = document.getElementById('prompt-input').value;
    let baseFolder = event.target.value;
    removeFolder(baseFolder, inputValue);
  }
  createPrompt('Enter the name of the folder you would like to delete.', confirmFolderName);
}

// This is the functioned called when Add Note button is clicked.
// This will take the click event info and feed it into the addNote function.
var clickAddNote = function (event) {
  let confirmNoteURL = function() {
    let inputValue = document.getElementById('prompt-input').value;
    let baseFolder = event.target.value;
    addNote(baseFolder, inputValue);
  }
  createPrompt('Copy and paste the URL of the note you would like to add.', confirmNoteURL);
}

// This is the function called when the remove note button is clicked.
// It will feed the event info to the removeNote function.
var clickRemoveNote = function (event) {
  let baseFolder = event.target.getAttribute('value');
  let folderNoteId = "" + event.target.id;
  let removeConfirm = function () {
    removeNote(baseFolder, folderNoteId);
  }
  createBoolPrompt("Are you sure you would like to remove this note?", removeConfirm);
};

var clickCopyNote = function (event) {
  let noteId = event.target.getAttribute('value');
  copyNote(noteId);
}

////////////////////  FUNCTIONS WITH API CALLS ////////////////////

// This is the function to open a folder using an api call.
// This will also open a folders menu and content and filling them. Additionally, this will clear and close all other menus.
var openFolder = function (folderID) {
  let idString = '' + folderID;
  $.post('/folder/api/openfolder', { FolderID: idString }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      let parentId = document.getElementById("" + folderID).getAttribute('value');

      // Set all open-menus to menus and clear them.
      let menus = document.getElementsByClassName('open-menu');
        for (let menuDiv of menus) {
          menuDiv.innerHTML = "";
          $(menuDiv).removeClass('open-menu');
          $(menuDiv).addClass('menu');
        }

      // Change the current folders menu to open and add its content.
      $(document.getElementById("" + folderID + '-menu')).append(
          '<button id="' + folderID + '-remove-folder" class="remove-folder" value="' + folderID + '"> Delete Folder </button>' +
          '<button id="' + folderID + '-rename-folder" class="rename-folder" value="' + folderID + '"> Rename Folder </button>' +
          '<button id="' + folderID + '-add-folder" class="add-folder" value="' + folderID + '"> New Folder </button>' +
          '<button id="' + folderID + '-add-note" class="add-note" value="' + folderID + '"> Add Note </button>');
      $(document.getElementById("" + folderID + '-menu')).removeClass('menu');
      $(document.getElementById("" + folderID + '-menu')).addClass('open-menu');
      $(document.getElementById("" + folderID + "-add-folder")).unbind();
      $(document.getElementById("" + folderID + "-add-folder")).click(clickAddFolder);
      $(document.getElementById("" + folderID + "-add-note")).unbind();
      $(document.getElementById("" + folderID + "-add-note")).click(clickAddNote);
      $(document.getElementById("" + folderID + "-remove-folder")).unbind()
      $(document.getElementById("" + folderID + "-remove-folder")).click(clickRemoveFolder);
      $(document.getElementById("" + folderID + "-rename-folder")).unbind()
      $(document.getElementById("" + folderID + "-rename-folder")).click(clickRenameFolder);


      // Append all child folder divs into the folder content.
      for (let referenceName of dataObj.folders) {
        if(referenceName != "") {
          let referenceId = "" + folderID + "/" + referenceName;
          $(document.getElementById("" + folderID + '-content')).append(
              '<div id="' + referenceId + '" class="folder" value="' + folderID + '"> ' + referenceName + ' </div>' +
              '<div id="' + referenceId + '-menu" class="menu"></div> ' +
              '<div id="' + referenceId + '-content" class="content"></div> ');
          $(document.getElementById(referenceId)).unbind();
          $(document.getElementById(referenceId)).click(clickFolder);
        }
      }

      //Append all child note divs into the folder content
      for (let note of dataObj.notes) {
        let noteId = "" + folderID + "/" + note.id;
        $(document.getElementById("" + folderID + '-content')).append(
            '<div id="' + noteId + '-container" class="note-container">' +
              '<div id="' + noteId + '-remove-note" class="remove-note" value="' + folderID + '"> X </div>' +
              '<div id="' + noteId + '-copy-note" class="copy-note" value="' + note.id + '"> Copy </div>' +
              '<div id="' + noteId + '" class="note" value="' + note.id + '""> ' + note.title + '</div>' +
            '</div>');
        $(document.getElementById(noteId + "-remove-note")).unbind();
        $(document.getElementById(noteId + "-remove-note")).click(clickRemoveNote);
        $(document.getElementById(noteId + "-copy-note")).unbind();
        $(document.getElementById(noteId + "-copy-note")).click(clickCopyNote);
        $(document.getElementById(noteId)).unbind();
        $(document.getElementById(noteId)).click(openNote);
      }
      $(document.getElementById('' + folderID + '-content')).removeClass("content");
      $(document.getElementById('' + folderID + '-content')).addClass("open-content");
      $(document.getElementById('' + folderID)).removeClass("folder");
      $(document.getElementById('' + folderID)).addClass("open-folder");
      if(document.getElementById('' + folderID + '-content').innerHTML == "") {
        document.getElementById('' + folderID + '-content').innerHTML = "(empty)";
      }
    }
    else {
      if (dataObj.code == 0) {
        displayError("An error has occured.", "Unable to open folder.");
      }
      if (dataObj.code == 1) {
        displayError("A database error has occured.");
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
    // Uh oh, something went wrong;
};

// This funcion makes an api call to add a folder with the given parameters.
var addFolder = function (parentId, folderName) {
  let parentIdString = "" + parentId;
  $.post('/folder/api/newfolder', { ParentID: parentIdString, FolderName: folderName }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      refreshContent(parentId);
    }
    else {
      if (dataObj.code == 0) {
        displayError("An error has occured.", "Ensure that a folder of the same name does not already exist in the parent folder.");
      }
      if (dataObj.code == 1) {
        displayError("A database error has occured.");
      }
      if (dataObj.code == 2) {
        displayError("The filepath for this folder is exceeds the 400 character limit.");
      }
      if (dataObj.code == 3) {
        displayError("You do not own the parent folder.");
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
};

var renameFolder = function (parentId, folderName, newFolderName) {
  let parentIdString = "" + parentId;
  $.post('/folder/api/renamefolder', { ParentID: parentIdString, FolderName: folderName, NewName: newFolderName }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      refreshContent(parentId);
    }
    else {
      if (dataObj.code == 0) {
        displayError("An error has occured.");
      }
      if (dataObj.code == 1) {
        displayError("A database error has occured.", "Ensure that the folder name was correctly entered.");
      }
      if (dataObj.code == 3) {
        displayError("You do not own the parent folder.");
      }
      if (dataObj.code == 4) {
        displayError("The folder is ROOT", "You cannot rename your base folder.");
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
};

// This funcion makes an api call to remove a folder with the given parameters.
var removeFolder = function (parentId, folderName) {
  let parentIdString = "" + parentId;
  $.post('/folder/api/deletefolder', { ParentID: parentIdString, FolderName: folderName }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      refreshContent(parentId);
    }
    else {
      if (dataObj.code == 0) {
        displayError("An error has occured.");
      }
      if (dataObj.code == 1) {
        displayError("A database error has occured.", "Ensure that the value entered matches the name of the folder you are trying to delete.");
      }
      if (dataObj.code == 3) {
        displayError("You do not own the parent folder.");
      }
      if (dataObj.code == 4) {
        displayError("The folder is ROOT", "You cannot delete your base folder.");
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
};

// This funcion makes an api call to add a note with the given parameters.
var addNote = function (parentId, noteId) {
  let parentIdString = "" + parentId;
  let noteIdInt = parseInt(noteId.substring(noteId.lastIndexOf('/') + 1), 10);
  $.post('/folder/api/addnote', { FolderID: parentIdString, NoteID: noteIdInt }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      refreshContent(parentId);
    }
    else {
      if (dataObj.code == 0) {
        displayError("An error has occured.", "Ensure that the given note is not already in this folder.");
      }
      if (dataObj.code == 1) {
        displayError("A database error has occured.", "Ensure that you have pasted a valid link into the input field.");
      }
      if (dataObj.code == 3) {
        displayError("You do not own the parent folder.");
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
};

// This funcion makes an api call to remove a folder with the given parameters.
var removeNote = function (parentId, folderNoteId) {
    let parentIdString = "" + parentId;
    let noteId = folderNoteId.substring(folderNoteId.lastIndexOf('/') + 1, (folderNoteId.indexOf('-')));
    let noteIdInt = parseInt(noteId, 10);
    $.post('/folder/api/removenote', { FolderID: parentIdString, NoteID: noteIdInt }, function (data) {
      let dataObj = $.parseJSON(data);
      if(dataObj.success == true) {
        refreshContent(parentId);
      }
      else {
        if (dataObj.code == 0) {
          displayError("An error has occured.");
        }
        if (dataObj.code == 1) {
          displayError("A database error has occured.");
        }
        if (dataObj.code == 3) {
          displayError("You do not own the parent folder.");
        }
      }
    }).fail(function () {
      displayError("Post request failed.");
    });
};

var copyNote = function (noteId) {
    let noteIdInt = parseInt(noteId, 10);
    console.log(noteIdInt);
    $.post('/note/api/copynote', { NoteID: noteIdInt }, function (data) {
      let dataObj = $.parseJSON(data);
      console.log(dataObj);
      if(dataObj.success == true) {
        console.log("note-copied");
      }
      else {
        if (dataObj.code == 0) {
          displayError("An error has occured.");
        }
        if (dataObj.code == 1) {
          displayError("A database error has occured.");
        }
      }
    }).fail(function () {
      displayError("Post request failed.");
    });
};



// This funcion makes an api call to initialize the root folder. It is called when the javascript is loaded.
var initializeRoot = function (rootId) {
  $.post('/folder/api/initializeroot', { RootID: rootId }, function (data) {
    let dataObj = $.parseJSON(data);

    if(!((dataObj.code == -1) || (dataObj.code == 6))) {
      if (dataObj.code == 0) {
        displayError("An error has occured when initializing the root folder.");
      }
      if (dataObj.code == 5) {
        displayError("The ID of the root folder does not match the ID of the user.");
      }
    }
    else {
      var rootArray = document.getElementsByClassName('root');
      for (let root of rootArray) {
        openFolder(root.id);
      }
    }
  }).fail(function () {
    displayError("Post request failed.");
  });
};

//////////////////// POPUP FUNCTIONS ////////////////////

// This function is assigned to a prompt-box's close button.
// It will hide the prompt box.
var closePrompt = function() {
  document.getElementById("prompt-box").style.visibility = "hidden";
};


// This function fills the promp-box with parameter info and makes the prompt box visible.
// This is used to get an input value from a user when required.
var createPrompt = function(message, confirmFunction) {
  document.getElementById('prompt-box').style.visibility = "visible";
  document.getElementById('prompt-box').innerHTML = '<div id="prompt-close">x</div>' +
      '<div id="prompt-message">' + message + '</div>' +
      '<input id="prompt-input" type="text" name="prompt-input" value=""/>' +
      '<button id="prompt-confirm"> Confirm </button>';
  $(document.getElementById('prompt-confirm')).unbind();
  $(document.getElementById('prompt-confirm')).click(confirmFunction);
  $(document.getElementById('prompt-confirm')).click(closePrompt);
  $(document.getElementById('prompt-close')).unbind(closePrompt);
  $(document.getElementById('prompt-close')).click(closePrompt);
};

var createDoublePrompt = function(message, confirmFunction) {
  document.getElementById('prompt-box').style.visibility = "visible";
  document.getElementById('prompt-box').innerHTML = '<div id="prompt-close">x</div>' +
      '<div id="prompt-message">' + message + '</div>' +
      '<input id="prompt-input" type="text" name="prompt-input" value=""/>' +
      '<input id="prompt-input2" type="text" name="prompt-input2" value=""/>' +
      '<button id="prompt-confirm"> Confirm </button>';
  $(document.getElementById('prompt-confirm')).unbind();
  $(document.getElementById('prompt-confirm')).click(confirmFunction);
  $(document.getElementById('prompt-confirm')).click(closePrompt);
  $(document.getElementById('prompt-close')).unbind(closePrompt);
  $(document.getElementById('prompt-close')).click(closePrompt);
};

// This function fills the prompt-box with parameter info and makes the prompt box visible.
// This is used when a user should confirm they intend to execute a function.
var createBoolPrompt = function (message, confirmFunction) {
  document.getElementById('prompt-box').style.visibility = "visible";
  document.getElementById('prompt-box').innerHTML = '<div id="prompt-close">x</div>' +
      '<div id="prompt-message">' + message + '</div>' +
    '<div style="text-align:center"> ' +
      '<button id="prompt-cancel"> Cancel </button>' +
      '<button id="prompt-confirm"> Confirm </button>' +
    '</div>';
  $(document.getElementById('prompt-confirm')).unbind();
  $(document.getElementById('prompt-confirm')).click(confirmFunction);
  $(document.getElementById('prompt-confirm')).click(closePrompt);
  $(document.getElementById('prompt-close')).unbind(closePrompt);
  $(document.getElementById('prompt-close')).click(closePrompt);
  $(document.getElementById('prompt-cancel')).unbind(closePrompt);
  $(document.getElementById('prompt-cancel')).click(closePrompt);
}

// This function is assigned to an error-box's close button.
// It will hide the error box.
var closeError = function() {
  document.getElementById("error-box").style.visibility = "hidden";
};

// This function fills the error-box with parameter info and makes the error box visible.
// This is used to alert a user that an error has occured.
var displayError = function(leadMessage, followMessage) {
  document.getElementById('error-box').style.visibility = "visible";
  if(followMessage) {
    document.getElementById('error-box').innerHTML = '<div id="error-close">x</div>' +
        '<div id="error-message"><strong> Error: ' + leadMessage + ' </strong> <br>' + followMessage + '</div>';
  }
  else {
    document.getElementById('error-box').innerHTML = '<div id="error-close">x</div>' +
        '<div id="error-message"><strong> Error: ' + leadMessage + ' </strong></div>';
  }
  $(document.getElementById('error-close')).click(closeError);
};

//////////////////// ONLOAD OPERATIONS ////////////////////

// This handles initializing the root folder and opening it when the .js file is initially loaded.
var rootArray = document.getElementsByClassName('root');
for (let root of rootArray) {
  initializeRoot(root.id);
  $(root).click(clickFolder);
  $(document.getElementById(root.id)).unbind();
  $(document.getElementById(root.id)).click(clickOpenFolder);
}


}
