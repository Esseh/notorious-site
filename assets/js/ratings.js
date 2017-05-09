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
      document.getElementById("current-rating").innerHTML = "Avg. Rating: " +
      dataObj.totalRating + " / 5";
    }
    else {
      console.log(dataObj.code);
    }
  }).fail(function () {
    console.log("Post request failed.");
  });
}

// 5
// 1-5

// 4
// 2-5

// 3
// 3-5

// 2
// 4-5


function noteRated(rating) {
  let i = 5 - rating + 1;
  let j = 5 - rating;
  while (i <= 5) {
    document.getElementById("star-" + i).style.backgroundImage = "url(/assets/images/FullStar.png)";
    i = i + 1;
  }
  while (j >= 1) {
    document.getElementById("star-" + j).style.backgroundImage = "url(/assets/images/EmptyStar.png)";
    j = j - 1;
  }
}

console.log(document.getElementById("current-rating"));
console.log(document.getElementById("current-rating").getAttribute('noteid'));
getRating(document.getElementById("current-rating").getAttribute('noteid'));