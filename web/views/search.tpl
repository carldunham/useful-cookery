{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/search.css">
{% endblock styles %}

{% block content %}
<div class="search">
  <div class="searchform">
    <form>
      <label for="query">Search: </label>
      <input type="search" id="query" name="q" size="40" placeholder="Enter Search Term(s)" value="{{ q }}">
      <input type="submit" value="Search">
    </form>
  </div>
  {% if error %}
    <div class="error">{{ errorText or 'Search error.' }}</div>
  {% endif %}
  {% if results %}
    <div class="summary">
      <span class="count">{{ results.stats.nfound }} result{% if results.stats.nfound != 1 %}s{% endif %} found{% if results.stats.nfound > 0 %} (in {{ "%0.3f"|format(results.stats.timeMicros/1000.0) }} seconds):{% endif %}</span>
    </div>
    <div class="results">
      {% for result in results.results %}
        <div class="result">
          <a href="/recipe/{{ result.obj.name }}">{{ result.obj.title }}</a> 
	  <span class="description">{{ result.obj.description }}</span>
	  <!-- <span class="score">{{ "%0.2f"|format(result.score) }}</span> -->
	</div>
      {% endfor %}
    </div>
  {% endif %}
</div>
{% endblock %}
