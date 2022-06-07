var isLoading = false
function postToLocal () { 
    var email = document.getElementById("email").innerText
    var token = document.getElementById("token").innerText
    var rToken = document.getElementById("rToken").innerText 
    var clientID = document.getElementById("clientID").innerText 
    var clientSecret = document.getElementById("clientSecret").innerText
    var port = document.getElementById("port").innerText
    
    document.getElementById("btn1").disabled= true
    isLoading = true
  
    url = 'http://localhost:' +  port + "/receiveToken"

    fetch(url, {
        method: "POST",
        headers: {"Content-Type":"application/json;charset=UTF-8"},
        body: JSON.stringify({ 
            token: token,
            email: email,
            rToken: rToken, 
            clientID : clientID, 
            clientSecret: clientSecret,
        }),
    })
    .then(data=> {
        isLoading = false
        document.getElementById("btn1").className = "w-48 inline-flex justify-center items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-gray-500 focus:outline-none bg-gray-200 hover:bg-grey-200"
        document.getElementById("btntext").innerText = "Saved"  
        //console.log(data)
        //window.close()
    })
    .catch(err => console.log(err))
}

