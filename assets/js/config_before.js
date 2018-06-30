function deleteRow(button) {
    var group = button.parentNode.parentNode;
    for (var i = 0; i < group.childNodes.length; i++) {
        var child = group.childNodes[i];
        if (child.type === 'hidden' && child.nodeName === 'INPUT' && strEndsWith(child.name, 'deleted')) {
            if (group.classList.contains('deleted')) {
                child.value = 'false';
                group.classList.remove('deleted');
            } else {
                child.value = 'true';
                group.classList.add('deleted');
            }
        }
    }
}

// Some browsers don't support string.endsWith function.
function strEndsWith(str, suffix) {
    return str.match(suffix+"$") == suffix;
}
