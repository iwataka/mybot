var twitterAccounts = new Bloodhound({
  datumTokenizer: Bloodhound.tokenizers.whitespace,
  queryTokenizer: Bloodhound.tokenizers.whitespace,
  remote: {
    url: '/twitter/users/search/?q=%QUERY',
    wildcard: '%QUERY'
  }
});

var typeaheadSettings = {
  name: 'twitter-accounts',
  display: 'screen_name',
  source: twitterAccounts,
  templates: {
    empty: [
      '<div class="empty-message"',
      'Type anything here',
      '</div>'
    ].join('\n'),
    suggestion: Handlebars.compile([
      '<div>',
      '<img src="{{profile_image_url}}" alt="profile image" height="20" width="20"/>',
      '<span>{{name}}@{{screen_name}}</span>',
      '<span class="label label-primary">followers: {{followers_count}}</span>',
      '</div>'
    ].join('\n'))
  }
}

$('input.typeahead-tagsinput').tagsinput({typeaheadjs: typeaheadSettings});

$(document).ready(function() {
  $('select.bootstrap-multiselect').multiselect({
    includeSelectAllOption: true,
    enableFiltering: true,
    enableCaseInsensitiveFiltering: true
  });
});
