var search_store = {
    keyword          : "",
    dicts            : dicts,
    relations        : relations,
    filtered_results : [],
    showing_results  : [],
    unshowing_results: [],
    count_per_page   : 32,
    is_loading       : false,
    loading_timer    : null,

    /**
     * init
     */

    init() {
        var url_params = new URLSearchParams(location.search);
        var q = url_params.get("q"); // keyword

        if (q !== null && q !== '') {
            this.keyword = q.trim();
        }
        this.search();
    },

    /**
     * updateParams
     */

    updateParams() {
        var searchParams = new URLSearchParams(window.location.search);

        if (this.keyword.trim() !== '') {
            searchParams.set("q", this.keyword.trim());
        } else {
            searchParams.delete("q")
        }
        window.history.replaceState(null, null, "?" + searchParams.toString());
    },

    /**
     * searchDatabase
     */

    searchDatabase() {
        this.filtered_results = this.dicts;
    },

    /**
     * search
     */

    search() {
        this.is_loading = true;
        this.updateParams();

        this.unshowing_results = [];
        this.showing_results   = [];

        if (this.loading_timer !== null) {
            clearTimeout(this.loading_timer);
        }
        this.loading_timer = setTimeout(() => {
            this.searchDatabase();

            if (this.keyword.trim() === '') {
                this.unshowing_results = this.filtered_results
            } else {
                var fuse = new Fuse(this.filtered_results, {
                    keys: ['c', 'e'], // c = Word, e = Examples
                })
                this.unshowing_results = fuse.search(this.keyword.trim()).map(v => v.item)
            }
            this.showMore(true);
            this.is_loading = false
        }, Math.random() * (300 - 200) + 300);
    },

    /**
     * showMore
     */

    showMore(ignore_loading) {
        if (ignore_loading === true) {
            this.showing_results = [...this.showing_results, ...this.unshowing_results.slice(0, this.count_per_page)]
            this.unshowing_results = this.unshowing_results.slice(this.count_per_page)
            return
        }
        this.is_loading = true;

        if (this.loading_timer !== null) {
            clearTimeout(this.loading_timer);
        }

        this.loading_timer = setTimeout(() => {
            this.showing_results = [...this.showing_results, ...this.unshowing_results.slice(0, this.count_per_page)]
            this.unshowing_results = this.unshowing_results.slice(this.count_per_page)
            this.is_loading = false
        }, Math.random() * (300 - 200) + 300);
    }
}