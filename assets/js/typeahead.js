var bestPictures = new Bloodhound({
  datumTokenizer: Bloodhound.tokenizers.whitespace,
  queryTokenizer: Bloodhound.tokenizers.whitespace,
  remote: {
    url: '../twitter/users/search/?q=%QUERY',
    wildcard: '%QUERY'
  }
});

$('#typeahead .typeahead').typeahead(null, {
  name: 'best-pictures',
  display: 'screen_name',
  source: bestPictures,
  templates: {
    empty: [
      '<div class="empty-message"',
      'Type anything here',
      '</div>'
    ].join('\n'),
    suggestion: Handlebars.compile('<div><img src="{{profile_image_url}}" alt="profile image" height="42" width="42"/>{{name}}@{{screen_name}}</div>')
  }
});
