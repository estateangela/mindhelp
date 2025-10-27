class CounselingCenter {
  final String id;
  final String name;
  final String address;
  final String phone;
  final bool onlineCounseling;
  final double? latitude;
  final double? longitude;
  final String createdAt;
  final String updatedAt;

  CounselingCenter({
    required this.id,
    required this.name,
    required this.address,
    required this.phone,
    required this.onlineCounseling,
    this.latitude,
    this.longitude,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CounselingCenter.fromJson(Map<String, dynamic> json) {
    return CounselingCenter(
      id: json['id'] as String,
      name: json['name'] as String,
      address: json['address'] as String,
      phone: json['phone'] as String,
      onlineCounseling: json['online_counseling'] as bool,
      latitude: (json['latitude'] == null)
          ? null
          : (json['latitude'] is num
              ? (json['latitude'] as num).toDouble()
              : double.tryParse(json['latitude'].toString())),
      longitude: (json['longitude'] == null)
          ? null
          : (json['longitude'] is num
              ? (json['longitude'] as num).toDouble()
              : double.tryParse(json['longitude'].toString())),
      createdAt: json['created_at'] as String,
      updatedAt: json['updated_at'] as String,
    );
  }
}
