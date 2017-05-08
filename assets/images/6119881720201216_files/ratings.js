function clickStar(rating, noteID) {
  $.post('/note/api/setrating', { NoteID: noteID, RatingValue: rating }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      noteRated(rating);
      getRating(noteID);
      console.log(dataObj);
    }
    else {
      console.log(dataObj.code);
    }
  }).fail(function () {
    console.log("Post request failed.");
  });
}

function getRating(noteID) {
  $.post('/note/api/getrating', { NoteID: noteID }, function (data) {
    let dataObj = $.parseJSON(data);
    if(dataObj.success == true) {
      console.log(dataObj);
      document.getElementById("current-rating").innerHTML = "Rating: " +
      dataObj.totalRating;
    }
    else {
      console.log(dataObj.code);
    }
  }).fail(function () {
    console.log("Post request failed.");
  });
}



function noteRated(rating) {

}

console.log(document.getElementById("current-rating"));
console.log(document.getElementById("current-rating").getAttribute('noteid'));
getRating(document.getElementById("current-rating").getAttribute('noteid'));