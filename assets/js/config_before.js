function deleteRow(button, modal) {
    var tr = button.parentNode.parentNode;
    if (tr.className === 'deleted') {
        for (var i = 0; i < tr.childNodes.length; i++) {
            var child = tr.childNodes[i]
            if (child.type === 'hidden' && child.nodeName === 'INPUT') {
                child.value = 'false';
                tr.className = '';
                button.innerHTML = 'Delete';
            }
        }
    } else {
        for (var i = 0; i < tr.childNodes.length; i++) {
            var child = tr.childNodes[i]
            if (child.type === 'hidden' && child.nodeName === 'INPUT') {
                child.value = 'true';
                tr.className = 'deleted';
                button.innerHTML = 'Revert';
            }
        }
    }
}

function openTimeline() {
    var screenName = document.getElementById('twitter-search-timeline-value').value;
    var modalBody = document.getElementById('twitter-search-timeline-modal-body');
    modalBody.innerHTML = '';
    twttr.widgets.createTimeline(
        {
            sourceType: 'profile',
            screenName: screenName
        },
        modalBody,
        {
            tweetLimit: 10
        }
    )
}

function openFavorites() {
    var screenName = document.getElementById('twitter-search-favorites-value').value;
    var modalBody = document.getElementById('twitter-search-favorites-modal-body');
    axios.get('../twitter/favorites/list?count=10&screen_name=' + screenName)
        .then(function (res) {
            createTweetsView(res.data, modalBody);
        })
        .catch(function (err) {
            console.log(err);
        });
}

function openSearch() {
    var query = document.getElementById('twitter-search-search-value').value;
    var modalBody = document.getElementById('twitter-search-search-modal-body');
    axios.get('../twitter/search?count=10&q=' + query)
        .then(function (res) {
            createTweetsView(res.data.statuses, modalBody);
        })
        .catch(function (err) {
            console.log(err);
        });
}

function createTweetsView(title, tweets, div) {
    div.innerHTML = ''
    for (i = 0; i < tweets.length; i++) {
        var id = tweets[i].id_str;
        twttr.widgets.createTweet(id, div);
    }
}
