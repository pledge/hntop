var Story = React.createClass({
	ticks: 0,

	tick: function() {
		this.ticks++;
		this.setState({ticksProg: this.ticks * 10});
		if(this.ticks >= 10) {
			this.loadStory();
			this.ticks = 0;
		}
	},

	loadStory: function() {
		$.get(this.props.source, function(result) {
			if(this.isMounted()) {
				this.setState({
					by: result.by,
					id: result.id,
					title: result.title,
					url: result.url
				});
			}
		}.bind(this));
	},

	getInitialState: function() {
		return {
			by: '',
			id: '',
			title: '',
			url: ''
		};
	},

	componentDidMount: function() {
		this.interval = setInterval(this.tick, 1000);
		this.loadStory();
	},

	componentWillUnmount: function() {
		clearInterval(this.interval);
	},

	render: function() {
		return (
			<div>
				<div className="progress progress-striped">
					<div className="progress-bar progress-bar-warning" style={{width: this.state.ticksProg + "%"}}></div>
				</div>
				<div>@<a href={"/user/" + this.state.by}>{this.state.by}</a> - <a href={this.state.url}>{this.state.title}</a></div>
			</div>
		);
	}
});

React.render(
	<Story source="/story" />, document.getElementById("story")
);

