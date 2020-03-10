package resources



func Manage() []byte{

	return []byte(manage)
}

var manage = `
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
	}
	.history-item{
		margin: 10px;
		padding: 10px;
		border-top: 1px solid #CCC;
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
</style>
<body>

<div>
	<div class="links">
		<a href="./">Incidents</a><br/>
	</div>
</div>
<h1>Callr Manage</h1>

<div id="ongoing" style="display: inline-block; width: 50%; vertical-align: top;" >
<div class="container">
	<h2>On Call</h2>
	<div id="oncall-outlet"></div>
</div>
</div
><div id="history" style="display: inline-block; width: 50%" >
<div class="container">
	<h2>People</h2>
	<div id="people-outlet"></div>
</div>
</div>
	


</body>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/lodash@4.17.15/lodash.min.js"></script>
<script type="application/javascript">

	
</script>

</html>
`
