

function deleteTodo(name){
    fetch(`/${name}`, {method: "Delete"}).then(res =>{
        if (res.status == 200){
            window.location.pathname = "/todos"
        }
    }).catch(err => {
        alert("An error occured while deleting the todo", err.message)
    })
}