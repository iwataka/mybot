var config;

function init() {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.open('GET', '/api/config/', false);
  xmlHttp.send(null);
  var res = xmlHttp.responseText;
  config = JSON.parse(res);
}

function loadConfig(elem) {
  init();
  root = document.getElementById(elem);

  timelineLabel = document.createElement('label');
  timelineLabel.innerHTML = 'Twitter Timelines';
  root.appendChild(timelineLabel);
  root.appendChild(timelineTable());

  favoriteLabel = document.createElement('label');
  favoriteLabel.innerHTML = 'Twitter Favorites';
  root.appendChild(favoriteLabel);
  root.appendChild(favoriteTable());

  searchLabel = document.createElement('label');
  searchLabel.innerHTML = 'Twitter Searches';
  root.appendChild(searchLabel);
  root.appendChild(searchTable());
}

function timelineTable() {
  table = document.createElement('table');
  table.className = 'table';
  tHead = document.createElement('thead');
  table.appendChild(tHead);

  headerRow = document.createElement('tr');
  tHead.appendChild(headerRow);
  headerTitles = ['screen_name', 'exclude_replies', 'include_trs', 'count'];
  for (var i = 0; i < headerTitles.length; i++) {
    data = document.createElement('th');
    data.innerHTML = headerTitles[i];
    headerRow.appendChild(data)
  }

  body = document.createElement('tbody');
  table.appendChild(body);
  timelines = config.Twitter.Timelines;
  for (var i = 0; i < timelines.length; i++) {
    timeline = timelines[i];
    row = document.createElement('tr');
    body.appendChild(row);
    screenName = document.createElement('td');
    screenName.innerHTML = timeline.ScreenName || timeline.ScreenNames.toString();
    row.appendChild(screenName);
    excludeReplies = document.createElement('td');
    excludeReplies.innerHTML = timeline.ExcludeReplies || '';
    row.appendChild(excludeReplies);
    includeRts = document.createElement('td');
    includeRts.innerHTML = timeline.IncludeRts || '';
    row.appendChild(includeRts);
    count = document.createElement('td');
    count.innerHTML = timeline.Count || '';
    row.appendChild(count);
  }

  return table;
}

function favoriteTable() {
  table = document.createElement('table');
  table.className = 'table';
  tHead = document.createElement('thead');
  table.appendChild(tHead);

  headerRow = document.createElement('tr');
  tHead.appendChild(headerRow);
  headerTitles = ['screen_name', 'count'];
  for (var i = 0; i < headerTitles.length; i++) {
    data = document.createElement('th');
    data.innerHTML = headerTitles[i];
    headerRow.appendChild(data)
  }

  body = document.createElement('tbody');
  table.appendChild(body);
  favorites = config.Twitter.Favorites;
  for (var i = 0; i < favorites.length; i++) {
    favorite = favorites[i];
    row = document.createElement('tr');
    body.appendChild(row);
    screenName = document.createElement('td');
    screenName.innerHTML = favorite.ScreenName || favorite.ScreenNames.toString();
    row.appendChild(screenName);
    count = document.createElement('td');
    count.innerHTML = favorite.Count || '';
    row.appendChild(count);
  }

  return table;
}

function searchTable() {
  table = document.createElement('table');
  table.className = 'table';
  tHead = document.createElement('thead');
  table.appendChild(tHead);

  headerRow = document.createElement('tr');
  tHead.appendChild(headerRow);
  headerTitles = ['query', 'result_type', 'count'];
  for (var i = 0; i < headerTitles.length; i++) {
    data = document.createElement('th');
    data.innerHTML = headerTitles[i];
    headerRow.appendChild(data)
  }

  body = document.createElement('tbody');
  table.appendChild(body);
  searches = config.Twitter.Searches;
  for (var i = 0; i < searches.length; i++) {
    search = searches[i];
    row = document.createElement('tr');
    body.appendChild(row);
    query = document.createElement('td');
    query.innerHTML = search.Query || search.Queries.toString();
    row.appendChild(query);
    resultType = document.createElement('td');
    resultType.innerHTML = search.ResultType || '';
    row.appendChild(resultType);
    count = document.createElement('td');
    count.innerHTML = search.Count || '';
    row.appendChild(count);
  }

  return table;
}
