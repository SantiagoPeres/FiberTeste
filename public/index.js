function removeFromDb(username){
    fetch(`/delete?username=${username}`, {method: "Delete"}).then(res =>{
        if (res.status == 200){
            window.location.pathname = "/"
        }
    })
 }

function updateDb(username){
	let input = document.getElementById(username)
    let newusername = input.value
    fetch(`/uptade?oldusername=${username}&newusername=${newusername}`, {method:"PUT"}).then(res => {
        if (res.status == 200){
            alert("Database uptade")
                window.location.pathname = "/"
        }
    })
}

