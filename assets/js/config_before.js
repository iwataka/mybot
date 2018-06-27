function deleteRow(button) {
    var tr = button.parentNode.parentNode;
    for (var i = 0; i < tr.childNodes.length; i++) {
        var child = tr.childNodes[i];
        if (child.type === 'hidden' && child.nodeName === 'INPUT' && strEndsWith(child.name, 'deleted')) {
            if (tr.className === 'deleted') {
                child.value = 'false';
                tr.className = '';
                button.innerHTML = 'Delete';
            } else {
                child.value = 'true';
                tr.className = 'deleted';
                button.innerHTML = 'Revert';
            }
        }
    }
}

// Some browsers don't support string.endsWith function.
function strEndsWith(str, suffix) {
    return str.match(suffix+"$") == suffix;
}
