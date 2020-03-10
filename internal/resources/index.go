package resources



func Index() []byte{

	return []byte(index)
}

var index = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Callr</title>
</head>
<style>
	body{
		font-family: Arial, sans-serif ;
	}
	h1{
		text-align: center;
	}
	a{
		cursor: pointer;
	}


	.container{
		margin: 10px;
		padding: 10px;
		border: 1px solid #ccc;
		overflow-x: hidden;
	}
	.history-item{
		margin: 10px;
		padding: 10px;
		border-top: 1px solid #CCC;
	}
	code{
		white-space: pre;
	}

	body .ctrl{
		display: none;
	}
	body.on-going .ctrl{
		display: inherit;
	}
	.log-entry{
		margin: 0 -10px;
		padding: 10px;
		background-color: #eee;
	}
	.log-entry:nth-child(even) {
	  background: #fff;
	}
	.timestamp{
		font-size: 0.9em;
	}

	.ctrl{
		float: right;
		padding: 20px;
		background: #e9e9e9;
		margin: -10px;
		border-left: 1px solid #CCC;
		border-bottom: 1px solid #CCC;
	}
	.ctrl button, button.view-log-button{
		cursor: pointer;
		width: 100%;
		background-color: #d5d5d5;
		color: #454545;
		font-weight: bold;
		border: 1px solid #7e7e7e;
		padding: 5px;
		margin: 4px;
		box-shadow: 1px 1px 2px 1px rgba(0, 0, 0, 0.2)
	}
	.ctrl button:hover, button.view-log-button:hover{
		background-color: #e4e4e4;
		box-shadow: 1px 1px 3px 1px rgba(0, 0, 0, 0.2)
	}
	button.view-log-button{
		width: initial;
		float: right;
	}
	.links{
		margin-top: -10px;
		float: right;
		padding: 10px;
		text-align: right;
		line-height: 1.4em;
	}
	.links a{
		color: gray;
		text-decoration: none;
	}
	.links a:hover{
		color: #454545;
		text-decoration: none;
	}

	#ongoing, #history{
		width: 50%;
		display: inline-block;
	}
	body.left #ongoing{
		width: 80%;
	}
	body.left #history{
		width: 20%;
	}
	body.right #ongoing{
		width: 20%;
	}
	body.right #history{
		width: 80%;
	}
</style>
<body class="left">
<div>
	<div class="links">
		<a href="manage">Manage</a><br/>
		<a onclick="toggleView()">View</a>
	</div>
</div>
<h1>Callr Incidents</h1>


<div id="ongoing" style="vertical-align: top;" >
<div class="container">
	<div class="ctrl">
		<button id="close-btn" class="btn" onclick="loadOngoing()">Refresh</button> <br/>
		<button id="close-btn" class="btn" onclick="closeIncident()">Close</button>
	</div>
	<h2>On-going		
	</h2>
	<div id="ongoing-outlet"></div>
</div>
</div
><div id="history" style="display: inline-block;" >
<div class="container">
	<h2>History</h2>
	<div id="history-outlet"></div>
</div>
</div>
	


</body>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/lodash@4.17.15/lodash.min.js"></script>
<script type="application/javascript">
	
	
	let HISTORY = [];
	let ONGOING = null;
	function closeIncident(){
	    if (!confirm("Are you sure you want to close the Incident?")){
	        return
	    }	    
	    console.log("closing");
	    $.ajax({
		  url: "incident",
		  method: "DELETE"
		}).done((res) =>{
		    console.log("Closed", res);
		    loadHistory();
		    loadOngoing();
		})
	}
	
	function toId(id){
	    return "id" + id.replace(/:/g, '')
	}
	
	function loadLog(id){
	    $.get("api/incidents/" + id + "/log").done((logs) => {
	        let dom = "";
	        for(let log of logs){
	            dom += parseLog(log)
	        }
	        $("#log" + toId(id)).html(dom)
	    });
	}
	function parseLog(log){
		if(log.content_type === 'application/json'){
			try{
				log.body = JSON.stringify(JSON.parse(log.body), null, 2)
			}catch (e) {}
		}
		
		if(log.content_type === 'application/x-www-form-urlencoded'){
		    try{
				log.body = JSON.stringify(query2Obj(log.body), null, 2)
			}catch (e) {}
		}
	    
	    let dom = '<div class="log-entry">';
		dom += '<div class="timestamp">' + new Date(log.created_at).toLocaleString('sv-SE') + '</div>';
		
		if(Object.keys(log.params).length !== 0 ) {
			dom +='<div><code><strong>Params</strong>\n';
			dom += JSON.stringify(log.params, null, 2) + '</code></div>';    
		}
		
		
		dom += '<div><code><strong>Content</strong>\n';
		dom += log.body + '</code></div>';
		
		dom += '</div>';
		return dom;
	}
	
	
	function loadHistory() {
	  $.get("api/incidents").done((res) => {
	      HISTORY = res;
	      dom = "";
	      for(let i of res){
	          dom += '<div class="history-item" id="'+ toId(i.id) +'"><a class="timestamp" onclick="expand(\'' + i.id + '\')">' + new Date(i.id).toLocaleString('sv-SE') + '</a></div>'
	      }
	      $("#history-outlet").html(dom);
	  })
	}
	
	function loadOngoing() {
	  $.get("api/incident").done((res) => {
	      ONGOING = res;
	      if(ONGOING){
	          let dom = '<div class="timestamp">' + new Date(ONGOING.created_at).toLocaleString('sv-SE') + '</div>';
	 		  dom += mkIncident(ONGOING);
	 		  console.log("dom", dom)
			  $("#ongoing-outlet").html(dom);
			  loadLog(ONGOING.id);
			  $("body").addClass("on-going");
	      }else{
	          $("#ongoing-outlet").html("");
	          $("body").removeClass("on-going");
	      }    
	  })
	}
	
	
	function mkIncident(i, expandable){
	    let dom = '';
		dom += '<p><code><strong>Incident:</strong><br/>';
		dom += JSON.stringify(i, null, 2);
		dom += '</code></p>';
		dom += '<p><strong>Log</strong>';
		if(expandable){
		    dom += '<button class="view-log-button" onclick="loadLog(\'' + i.id + '\')">View</button>'
		}
	 	dom +='</p>';
		dom += '<p id="log' + toId(i.id) +'">';		
		dom += '</p>'
		return dom
	}
	
	
	function minimize(id) {
	  $("#" +toId(id)).html('<a class="timestamp" onclick="expand(\'' + id + '\')">' +  new Date(id).toLocaleString('sv-SE')  + '</a>')
	}
	
	function expand(id) {
	   let i =  _.find(HISTORY, {'id': id});
	  
	   let dom = '<a class="timestamp" onclick="minimize(\'' + i.id + '\')">' +  new Date(i.id).toLocaleString('sv-SE')  + '</a>';
	   dom += mkIncident(i, true);
	  
	   console.log(dom);
	   $("#" + toId(id)).html(dom)
	}	
	
	function query2Obj(str) {
	    let res = {};
	    let parts = str.split('&');
	    for(let part of parts){
	        if(part.length === 0){
	            continue;
	        }
	        let pp = part.split("=");
	        res[pp[0]] = decodeURIComponent(pp[1]);
	    } 
	  	return res
	}

	let VIEW = 0; 
	function toggleView() {
	  VIEW++;
	  
	  let b = $("body")
	  b.removeClass("left");
	  b.removeClass("right");
	  
	  switch (VIEW%3) {
	      case 0:
	        b.addClass("left");
		   	break;
		  case 1: break;
		  case 2:
		   b.addClass("right");
		    break;
	  }  
	}
	
	loadOngoing();
	loadHistory();
</script>

</html>
`
